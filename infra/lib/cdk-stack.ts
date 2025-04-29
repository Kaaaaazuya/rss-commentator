import * as cdk from "aws-cdk-lib";
import * as events from "aws-cdk-lib/aws-events";
import * as ecr from "aws-cdk-lib/aws-ecr";
import * as dynamodb from "aws-cdk-lib/aws-dynamodb";
import * as iam from "aws-cdk-lib/aws-iam";
import type { Construct } from "constructs";
import { Serverless } from "./constructs/serverless";

export interface RssServiceStackProps extends cdk.StackProps {
	/**
	 * 既存のリソースを再利用するかどうか
	 * trueの場合、存在していればそれを使用し、存在しなければ作成します
	 */
	readonly reuseExistingResource?: boolean;
}

export class RssServiceStack extends cdk.Stack {
	public readonly articlesTable: dynamodb.Table;
	public readonly tagsTable: dynamodb.Table;
	public readonly articlesTagsTable: dynamodb.Table;
	public readonly lambdaRole: iam.Role;

	constructor(scope: Construct, id: string, props?: RssServiceStackProps) {
		super(scope, id, props);

		const reuseExistingResource =
			props?.reuseExistingResource ||
			process.env.REUSE_EXISTING_RESOURCE === "true";

		// ------ DynamoDB -----
		// テーブル作成またはインポート
		this.articlesTable = this.createOrImportTable(
			"ArticlesTable",
			"articles",
			{ name: "url_hash", type: dynamodb.AttributeType.STRING },
			undefined,
			reuseExistingResource,
		);

		this.tagsTable = this.createOrImportTable(
			"TagsTable",
			"tags",
			{ name: "tag_name", type: dynamodb.AttributeType.STRING },
			undefined,
			reuseExistingResource,
		);

		this.articlesTagsTable = this.createOrImportTable(
			"ArticlesTagsTable",
			"articles_tags",
			{ name: "url_hash", type: dynamodb.AttributeType.STRING },
			{ name: "tag_name", type: dynamodb.AttributeType.STRING },
			reuseExistingResource,
		);

		// ECR リポジトリの作成
		const rssFetcherRepo = this.createEcrRepository(
			"RssFetcherRepo",
			"rss-fetcher",
			cdk.RemovalPolicy.DESTROY,
			reuseExistingResource,
		);
		const summarizerRepo = this.createEcrRepository(
			"SummarizerRepo",
			"summarizer",
			cdk.RemovalPolicy.DESTROY,
			reuseExistingResource,
		);
		// notifierリポジトリ名を修正
		const summaryNotifierRepo = this.createEcrRepository(
			"SummaryNotifierRepo",
			"notifier",
			cdk.RemovalPolicy.DESTROY,
			reuseExistingResource,
		);

		// ----- Lambda -----
		// fetcher: 毎日21時に起動
		const fetcherRule = new events.Rule(this, "FetchRule", {
			schedule: events.Schedule.cron({ minute: "0", hour: "21" }),
		});
		const fetcher = new Serverless(this, "RSSFetcher", {
			functionName: "rss-fetcher",
			functionImageRepositoryName: rssFetcherRepo.repositoryName,
			functionIamRolePolicies: [
				new iam.PolicyStatement({
					effect: iam.Effect.ALLOW,
					actions: ["dynamodb:PutItem"],
					resources: [this.articlesTable.tableArn],
				}),
			],
			eventSource: fetcherRule,
		});

		// summarizer: 毎日21時20分に起動
		const summarizerRule = new events.Rule(this, "SummarizeRule", {
			schedule: events.Schedule.cron({ minute: "20", hour: "21" }),
		});
		const summarizer = new Serverless(this, "Summarizer", {
			functionName: "summarizer",
			functionImageRepositoryName: summarizerRepo.repositoryName,
			functionIamRolePolicies: [
				new iam.PolicyStatement({
					effect: iam.Effect.ALLOW,
					actions: ["dynamodb:PutItem", "dynamodb:Update*", "dynamodb:Scan"],
					resources: [
						this.articlesTable.tableArn,
						this.tagsTable.tableArn,
						this.articlesTagsTable.tableArn,
					],
				}),
			],
			eventSource: summarizerRule,
		});

		// notifier: 毎日21時40分に起動
		const notifierRule = new events.Rule(this, "NotifyRule", {
			schedule: events.Schedule.cron({ minute: "40", hour: "21" }),
		});
		const notifier = new Serverless(this, "Notifier", {
			functionName: "notifier",
			functionImageRepositoryName: summaryNotifierRepo.repositoryName,
			functionIamRolePolicies: [
				new iam.PolicyStatement({
					effect: iam.Effect.ALLOW,
					actions: ["dynamodb:Scan", "dynamodb:Query"],
					resources: [
						this.articlesTable.tableArn,
						this.articlesTagsTable.tableArn,
					],
				}),
			],
			eventSource: notifierRule,
		});

		// スタックからの出力を追加
		// DynamoDBテーブル
		new cdk.CfnOutput(this, "ArticlesTableName", {
			value: this.articlesTable.tableName,
			description: "Articles DynamoDB table name",
			exportName: "ArticlesTableName",
		});
		new cdk.CfnOutput(this, "TagsTableName", {
			value: this.tagsTable.tableName,
			description: "Tags DynamoDB table name",
			exportName: "TagsTableName",
		});
		new cdk.CfnOutput(this, "ArticlesTagsTableName", {
			value: this.articlesTagsTable.tableName,
			description: "Articles-Tags DynamoDB table name",
			exportName: "ArticlesTagsTableName",
		});

		// ECRリポジトリ
		new cdk.CfnOutput(this, "RssFetcherRepoName", {
			value: rssFetcherRepo.repositoryName,
			description: "RSS Fetcher ECR repository name",
			exportName: "RssFetcherRepoName",
		});
		new cdk.CfnOutput(this, "SummarizerRepoName", {
			value: summarizerRepo.repositoryName,
			description: "Summarizer ECR repository name",
			exportName: "SummarizerRepoName",
		});
		new cdk.CfnOutput(this, "SummaryNotifierRepoName", {
			value: summaryNotifierRepo.repositoryName,
			description: "Summary Notifier ECR repository name",
			exportName: "SummaryNotifierRepoName",
		});

		// Lambda関数
		new cdk.CfnOutput(this, "RssFetcherFunctionName", {
			value: fetcher.lambdaFunction.functionName,
			description: "RSS Fetcher Lambda function name",
			exportName: "RssFetcherFunctionName",
		});
		new cdk.CfnOutput(this, "SummarizerFunctionName", {
			value: summarizer.lambdaFunction.functionName,
			description: "Summarizer Lambda function name",
			exportName: "SummarizerFunctionName",
		});
		new cdk.CfnOutput(this, "NotifierFunctionName", {
			value: notifier.lambdaFunction.functionName,
			description: "Notifier Lambda function name",
			exportName: "NotifierFunctionName",
		});
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
		reuseExisting = false,
	): ecr.IRepository {
		if (reuseExisting) {
			try {
				return ecr.Repository.fromRepositoryName(this, id, repositoryName);
			} catch (error) {
				console.log(
					`ECRリポジトリ ${repositoryName} が存在しないため、新規作成します`,
				);
			}
		}
		return new ecr.Repository(this, id, {
			repositoryName,
			imageScanOnPush: true,
			emptyOnDelete: true,
			removalPolicy,
		});
	}

	/**
	 * テーブルを作成またはインポートするヘルパーメソッド
	 * @param id Construct内での一意なID
	 * @param tableName テーブル名
	 * @param partitionKey パーティションキー
	 * @param sortKey ソートキー（オプション）
	 * @param reuseExisting 既存のテーブルを再利用するかどうか
	 * @returns DynamoDBテーブル
	 */
	private createOrImportTable(
		id: string,
		tableName: string,
		partitionKey: { name: string; type: dynamodb.AttributeType },
		sortKey?: { name: string; type: dynamodb.AttributeType },
		reuseExisting = false,
	): dynamodb.Table {
		// 既存のテーブルを再利用する場合
		if (reuseExisting) {
			try {
				// 既存のテーブルをインポートしてみる
				return dynamodb.Table.fromTableName(
					this,
					id,
					tableName,
				) as dynamodb.Table;
			} catch (error) {
				console.log(`テーブル ${tableName} が存在しないため、新規作成します`);
				// インポートに失敗した場合は新規作成
			}
		}

		// 新しいテーブルを作成
		return new dynamodb.Table(this, id, {
			tableName,
			partitionKey,
			sortKey,
			billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
			removalPolicy: reuseExisting
				? cdk.RemovalPolicy.RETAIN
				: cdk.RemovalPolicy.DESTROY,
		});
	}
}
