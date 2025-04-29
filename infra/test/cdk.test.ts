import * as cdk from 'aws-cdk-lib';
import { Template, Match } from 'aws-cdk-lib/assertions';
import { RssServiceStack } from '../lib/cdk-stack';
import type { App } from 'aws-cdk-lib';

describe('RssServiceStack Resources', () => {
  let app: App;
  let stack: RssServiceStack;
  let template: Template;

  beforeAll(() => {
    app = new cdk.App();
    stack = new RssServiceStack(app, 'TestRssServiceStack');
    template = Template.fromStack(stack);
  });

  test('DynamoDB Table is created', () => {
    template.hasResourceProperties('AWS::DynamoDB::Table', {
      TableName: 'articles',
      BillingMode: 'PAY_PER_REQUEST',
      KeySchema: [{ AttributeName: 'canonical_url', KeyType: 'HASH' }],
      AttributeDefinitions: [{ AttributeName: 'canonical_url', AttributeType: 'S' }],
    });
  });

  test('RssFetcher Lambda is created with correct properties', () => {
    template.hasResourceProperties('AWS::Lambda::Function', {
      FunctionName: 'rss_fetcher',
      PackageType: 'Image',
      Architectures: ['arm64'],
      MemorySize: 512,
      Timeout: 30,
      Environment: {
        Variables: {
          // 生成されたトークンが含まれていることを検証
          ARTICLES_TABLE_NAME: {
            Ref: Match.stringLikeRegexp('ArticlesTable'),
          },
        },
      },
    });

    template.hasResourceProperties('AWS::Lambda::Url', {
      AuthType: 'NONE',
    });
  });

  test('ECR Repository is created', () => {
    template.hasResourceProperties('AWS::ECR::Repository', {
      RepositoryName: 'rss-commentator',
      ImageScanningConfiguration: {
        ScanOnPush: true,
      },
    });
  });

  test('IAM Role has correct permissions', () => {
    // Lambda 用 IAM ロールのトラストポリシーを検証
    template.hasResourceProperties('AWS::IAM::Role', {
      AssumeRolePolicyDocument: {
        Statement: [
          {
            Action: 'sts:AssumeRole',
            Effect: 'Allow',
            Principal: {
              Service: 'lambda.amazonaws.com',
            },
          },
        ],
      },
    });

    // grantReadWriteData により付与された DynamoDB アクションが含まれているか検証
    template.hasResourceProperties('AWS::IAM::Policy', {
      PolicyDocument: {
        Statement: Match.arrayWith([
          Match.objectLike({
            Effect: 'Allow',
            Action: Match.arrayWith([
              'dynamodb:PutItem',
              'dynamodb:GetItem',
              'dynamodb:UpdateItem',
              'dynamodb:DeleteItem',
              'dynamodb:Scan',
              'dynamodb:Query',
            ]),
            Resource: Match.anyValue(),
          }),
        ]),
      },
    });
  });
});
