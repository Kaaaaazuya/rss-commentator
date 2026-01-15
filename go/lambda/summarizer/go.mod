module github.com/Kaaaaazuya/rss-commentator/go/lambda/summarizer

go 1.24.4

toolchain go1.25.6

require (
	github.com/Kaaaaazuya/rss-commentator/go/shared v0.0.0
	github.com/aws/aws-lambda-go v1.51.2
	github.com/aws/aws-sdk-go-v2 v1.41.1
	github.com/aws/aws-sdk-go-v2/config v1.32.7
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.53.6
	github.com/tmc/langchaingo v0.1.14
	go.uber.org/zap v1.27.1
	gopkg.in/guregu/null.v3 v3.5.0
)

require (
	github.com/aws/aws-sdk-go-v2/credentials v1.19.7 // indirect
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.20.30 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.18.17 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.17 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.17 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodbstreams v1.32.10 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.11.17 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.17 // indirect
	github.com/aws/aws-sdk-go-v2/service/signin v1.0.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.30.9 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.35.13 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.41.6 // indirect
	github.com/aws/smithy-go v1.24.0 // indirect
	github.com/dlclark/regexp2 v1.10.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/pkoukk/tiktoken-go v0.1.6 // indirect
	go.uber.org/multierr v1.11.0 // indirect
)

replace github.com/Kaaaaazuya/rss-commentator/go/shared v0.0.0 => ../../shared
