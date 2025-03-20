import * as cdk from "aws-cdk-lib";
import * as lambda from "aws-cdk-lib/aws-lambda";
import * as ecr from "aws-cdk-lib/aws-ecr";
import * as dynamodb from "aws-cdk-lib/aws-dynamodb";
import * as iam from "aws-cdk-lib/aws-iam";
import type { Construct } from "constructs";

export class RssServiceStack extends cdk.Stack {
	public readonly articlesTable: dynamodb.Table;
	public readonly lambdaRole: iam.Role;

	constructor(scope: Construct, id: string, props?: cdk.StackProps) {
		super(scope, id, props);

		// DynamoDB Table の作成
		this.articlesTable = new dynamodb.Table(this, "ArticlesTable", {
			tableName: "articles",
			partitionKey: {
				name: "canonical_url",
				type: dynamodb.AttributeType.STRING,
			},
			billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
		});

		// Lambda 用 IAM ロールの作成
		this.lambdaRole = new iam.Role(this, "LambdaExecutionRole", {
			assumedBy: new iam.ServicePrincipal("lambda.amazonaws.com"),
			managedPolicies: [
				iam.ManagedPolicy.fromAwsManagedPolicyName(
					"service-role/AWSLambdaBasicExecutionRole",
				),
			],
		});
		// DynamoDB へのアクセス権を付与
		this.articlesTable.grantReadWriteData(this.lambdaRole);

		// ECR リポジトリの作成
		const rssFetcherRepo = this.createEcrRepository(
			"RssFetcherRepo",
			"rss-commentator",
		);

		// Docker イメージ Lambda の作成（rssFetcher）
		const rssFetcher = this.createDockerLambda({
			id: "RssFetcher",
			functionName: "rss_fetcher",
			ecrRepository: rssFetcherRepo,
			cmd: ["main"],
			environment: {
				ARTICLES_TABLE_NAME: this.articlesTable.tableName,
			},
			timeoutSeconds: 30,
			memorySize: 512,
			architecture: lambda.Architecture.ARM_64,
		});

		// 他の Lambda や ECR リポジトリを追加する際は同様に createDockerLambda や createEcrRepository を利用できます
	}

	/**
	 * ECR リポジトリ作成のためのヘルパーメソッド
	 * @param id Construct 内での一意な ID
	 * @param repositoryName リポジトリ名
	 * @param removalPolicy 削除ポリシー（デフォルトは DESTROY）
	 * @returns 作成された ECR Repository
	 */
	private createEcrRepository(
		id: string,
		repositoryName: string,
		removalPolicy: cdk.RemovalPolicy = cdk.RemovalPolicy.DESTROY,
	): ecr.Repository {
		return new ecr.Repository(this, id, {
			repositoryName,
			imageScanOnPush: true,
			removalPolicy,
		});
	}

	/**
	 * Docker イメージを利用した Lambda 関数作成のためのヘルパーメソッド
	 * @param params Docker Lambda 作成の各種パラメータ
	 * @returns 作成された DockerImageFunction
	 */
	private createDockerLambda(params: {
		id: string;
		functionName: string;
		ecrRepository: ecr.Repository;
		cmd: string[];
		environment?: { [key: string]: string };
		timeoutSeconds?: number;
		memorySize?: number;
		architecture?: lambda.Architecture;
	}): lambda.DockerImageFunction {
		const func = new lambda.DockerImageFunction(this, params.id, {
			functionName: params.functionName,
			code: lambda.DockerImageCode.fromEcr(params.ecrRepository, {
				cmd: params.cmd,
			}),
			role: this.lambdaRole,
			environment: params.environment,
			timeout: cdk.Duration.seconds(params.timeoutSeconds ?? 30),
			memorySize: params.memorySize ?? 512,
			architecture: params.architecture ?? lambda.Architecture.X86_64,
		});

		// Lambda Function URL の設定（必要に応じて）
		func.addFunctionUrl({
			authType: lambda.FunctionUrlAuthType.NONE,
		});

		return func;
	}
}
