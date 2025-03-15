import * as cdk from "aws-cdk-lib";
import * as lambda from "aws-cdk-lib/aws-lambda";
import * as dynamodb from "aws-cdk-lib/aws-dynamodb";
import * as iam from "aws-cdk-lib/aws-iam";
import type { Construct } from "constructs";

export class RssServiceStack extends cdk.Stack {
	public readonly articlesTable: dynamodb.Table;
	public readonly userApi: lambda.Function;

	constructor(scope: Construct, id: string, props?: cdk.StackProps) {
		super(scope, id, props);

		// DynamoDB Table
		this.articlesTable = new dynamodb.Table(this, "ArticlesTable", {
			tableName: "articles",
			partitionKey: {
				name: "canonical_url",
				type: dynamodb.AttributeType.STRING,
			},
			billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
		});

		// IAM Role for Lambda
		const lambdaRole = new iam.Role(this, "LambdaExecutionRole", {
			assumedBy: new iam.ServicePrincipal("lambda.amazonaws.com"),
			managedPolicies: [
				iam.ManagedPolicy.fromAwsManagedPolicyName(
					"service-role/AWSLambdaBasicExecutionRole",
				),
			],
		});

		// Allow Lambda to access DynamoDB
		this.articlesTable.grantReadWriteData(lambdaRole);

		// Define Lambda Functions
		const createLambda = (name: string) => {
			const fn = new lambda.Function(this, name, {
				functionName: name,
				runtime: lambda.Runtime.NODEJS_22_X,
				handler: "index.handler",
				code: lambda.Code.fromInline("exports.handler = async () => {};"), // Jest 実行時はインラインコードを使用
				role: lambdaRole,
				environment: {
					ARTICLES_TABLE_NAME: this.articlesTable.tableName,
				},
			});

			// Lambda Function URL
			fn.addFunctionUrl({
				authType: lambda.FunctionUrlAuthType.NONE, // 認証なし（必要に応じて変更）
			});

			return fn;
		};

		const rssFetcher = createLambda("rss_fetcher");
		const articleFetcher = createLambda("article_fetcher");
		const summarizer = createLambda("summarizer");
		const notifier = createLambda("notifier");
		this.userApi = createLambda("user_api");
	}
}
