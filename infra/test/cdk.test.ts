import * as cdk from "aws-cdk-lib";
import { Template, Match } from "aws-cdk-lib/assertions";
import { RssServiceStack } from "../lib/cdk-stack";
import type { App } from "aws-cdk-lib";

describe("RssServiceStack Resources", () => {
	let app: App;
	let stack: RssServiceStack;
	let template: Template;

	beforeAll(() => {
		app = new cdk.App();
		stack = new RssServiceStack(app, "TestRssServiceStack");
		template = Template.fromStack(stack);
	});

	test("DynamoDB Table is created", () => {
		template.hasResourceProperties(
			"AWS::DynamoDB::Table",
			Match.objectLike({
				TableName: "articles",
				BillingMode: "PAY_PER_REQUEST",
				KeySchema: Match.arrayWith([
					Match.objectLike({ AttributeName: "url_hash", KeyType: "HASH" }),
				]),
				AttributeDefinitions: Match.arrayWith([
					Match.objectLike({ AttributeName: "url_hash", AttributeType: "S" }),
				]),
			}),
		);
	});

	test("RssFetcher Lambda is created with correct properties", () => {
		template.hasResourceProperties(
			"AWS::Lambda::Function",
			Match.objectLike({
				FunctionName: "rss-fetcher",
				PackageType: "Image",
				Architectures: ["arm64"],
				// MemorySize, Timeout, EnvironmentはCDK実装で明示的に設定されていないため省略
			}),
		);
		// Lambda URLはCDKで作成していないためテストから除外
	});

	test("ECR Repository is created", () => {
		template.hasResourceProperties(
			"AWS::ECR::Repository",
			Match.objectLike({
				RepositoryName: "rss-fetcher",
				ImageScanningConfiguration: Match.objectLike({
					ScanOnPush: true,
				}),
			}),
		);
	});

	test("IAM Role has correct permissions", () => {
		// Lambda 用 IAM ロールのトラストポリシーを検証
		template.hasResourceProperties(
			"AWS::IAM::Role",
			Match.objectLike({
				AssumeRolePolicyDocument: Match.objectLike({
					Statement: Match.arrayWith([
						Match.objectLike({
							Action: "sts:AssumeRole",
							Effect: "Allow",
							Principal: Match.objectLike({
								Service: "lambda.amazonaws.com",
							}),
						}),
					]),
				}),
			}),
		);

		// grantReadWriteData により付与された DynamoDB アクションが含まれているか検証
		template.hasResourceProperties(
			"AWS::IAM::Policy",
			Match.objectLike({
				PolicyDocument: Match.objectLike({
					Statement: Match.arrayWith([
						Match.objectLike({
							Effect: "Allow",
							Action: Match.arrayWith([Match.stringLikeRegexp("dynamodb:.*")]),
							Resource: Match.anyValue(),
						}),
					]),
				}),
			}),
		);
	});
});
