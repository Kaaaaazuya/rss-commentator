import * as cdk from "aws-cdk-lib";
import { Template } from "aws-cdk-lib/assertions";
import { RssServiceStack } from "../lib/cdk-stack";
import type { App } from "aws-cdk-lib";

describe("RssServiceStack", () => {
	let app: App;
	let stack: RssServiceStack;
	let template: Template;

	beforeAll(() => {
		app = new cdk.App();
		stack = new RssServiceStack(app, "TestRssServiceStack");
		template = Template.fromStack(stack);
	});

	test("DynamoDB Table is created", () => {
		template.hasResourceProperties("AWS::DynamoDB::Table", {
			TableName: "articles",
			BillingMode: "PAY_PER_REQUEST",
			KeySchema: [{ AttributeName: "canonical_url", KeyType: "HASH" }],
			AttributeDefinitions: [
				{ AttributeName: "canonical_url", AttributeType: "S" },
			],
		});
	});

	test("Lambda Functions are created with Function URLs", () => {
		const lambdaNames = [
			"rss_fetcher",
			"article_fetcher",
			"summarizer",
			"notifier",
			"user_api",
		];

		// biome-ignore lint/complexity/noForEach: <explanation>
		lambdaNames.forEach((lambdaName) => {
			template.hasResourceProperties("AWS::Lambda::Function", {
				FunctionName: lambdaName,
				Runtime: "nodejs22.x",
				Handler: "index.handler",
			});

			template.hasResourceProperties("AWS::Lambda::Url", {
				AuthType: "NONE",
			});
		});
	});

	test("Snapshot Test", () => {
		expect(template.toJSON()).toMatchSnapshot();
	});
});
