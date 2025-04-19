import { Rule } from "aws-cdk-lib/aws-events";
import {
	LambdaFunction,
	LambdaFunctionProps,
} from "aws-cdk-lib/aws-events-targets";
import * as iam from "aws-cdk-lib/aws-iam";
import type * as ec2 from "aws-cdk-lib/aws-ec2";
import * as lambda from "aws-cdk-lib/aws-lambda";
import type { Queue } from "aws-cdk-lib/aws-sqs";
import { Construct } from "constructs";
import { Repository } from "aws-cdk-lib/aws-ecr";
import { RetentionDays } from "aws-cdk-lib/aws-logs";

/*
Lambda リソース定義
    * IAM ロール
    * コンテナイメージの Lambda 関数
    * イベントソース(EventBridge)
    * パラメータストア
デプロイは lambroll で管理
    * メモリやタイムアウトも lambroll の設定ファイルで管理する
 */

type ServerlessProps = {
	functionName: string;
	functionImageRepositoryName: string;
	functionIamRolePolicies: iam.PolicyStatement[];
	eventSource?: Queue | Rule;
	eventSourceProps?: LambdaFunctionProps;
};

export class Serverless extends Construct {
	public readonly securityGroup: ec2.SecurityGroup;
	public readonly lambdaFunction: lambda.DockerImageFunction;
	public readonly executionRole: iam.Role;

	constructor(scope: Construct, id: string, props: ServerlessProps) {
		super(scope, id);

		const repositoryName = props.functionImageRepositoryName;
		const repository = Repository.fromRepositoryName(
			this,
			"ImageRepository",
			repositoryName,
		);
		const functionName = props.functionName;

		// IAM Role
		this.executionRole = new iam.Role(this, "LambdaExecutionRole", {
			assumedBy: new iam.ServicePrincipal("lambda.amazonaws.com"),
			description: `${functionName} Lambda execution role`,
			managedPolicies: [
				iam.ManagedPolicy.fromAwsManagedPolicyName(
					"service-role/AWSLambdaBasicExecutionRole",
				),
			],
		});
		this.executionRole.addToPolicy(
			new iam.PolicyStatement({
				effect: iam.Effect.ALLOW,
				actions: ["ecr:SetRepositoryPolicy", "ecr:GetRepositoryPolicy"],
				resources: [repository.repositoryArn],
			}),
		);

		// カスタムポリシーを追加
		props.functionIamRolePolicies.forEach((policy) => {
			this.executionRole.addToPolicy(policy);
		});

		this.lambdaFunction = new lambda.DockerImageFunction(this, "Function", {
			functionName: functionName,
			code: lambda.DockerImageCode.fromEcr(repository, {
				tagOrDigest: "latest",
			}),
			logRetention: RetentionDays.ONE_WEEK,
			architecture: lambda.Architecture.ARM_64,
			role: this.executionRole,
		});
		if (props.eventSource instanceof Rule) {
			const ruleProps = props.eventSourceProps as
				| LambdaFunctionProps
				| undefined;
			props.eventSource.addTarget(
				new LambdaFunction(this.lambdaFunction, ruleProps),
			);
		}
	}
}
