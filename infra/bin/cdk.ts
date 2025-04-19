#!/usr/bin/env node
import * as cdk from "aws-cdk-lib";
import { RssServiceStack } from "../lib/cdk-stack";

const app = new cdk.App();
new RssServiceStack(app, "RssCommentator", {
	/* AWSアカウントとリージョンを環境変数から取得するように設定 */
	env: {
		account: process.env.CDK_DEFAULT_ACCOUNT,
		region: process.env.CDK_DEFAULT_REGION || "ap-northeast-1",
	},
	/* スタックに共通のタグを追加 */
	tags: {
		Project: "RssCommentator",
		Environment: process.env.ENVIRONMENT || "dev",
		ManagedBy: "CDK",
	},
});
