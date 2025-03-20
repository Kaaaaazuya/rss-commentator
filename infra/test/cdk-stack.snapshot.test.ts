import * as cdk from "aws-cdk-lib";
import { Template } from "aws-cdk-lib/assertions";
import { RssServiceStack } from "../lib/cdk-stack";
import type { App } from "aws-cdk-lib";

describe("RssServiceStack Snapshot", () => {
	let app: App;
	let stack: RssServiceStack;
	let template: Template;

	beforeAll(() => {
		app = new cdk.App();
		// スナップショット用に別のスタックIDを使用
		stack = new RssServiceStack(app, "TestRssServiceStackSnapshot");
		template = Template.fromStack(stack);
	});

	test("Snapshot Test", () => {
		expect(template.toJSON()).toMatchSnapshot();
	});
});
