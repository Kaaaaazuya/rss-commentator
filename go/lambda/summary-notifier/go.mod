module github.com/Kaaaaazuya/rss-commentator/go/lambda/summary-notifier

go 1.24.1

require (
	github.com/Kaaaaazuya/rss-commentator/go/shared v0.0.0
	github.com/aws/aws-lambda-go v1.54.0
	github.com/aws/aws-sdk-go-v2 v1.43.8
	github.com/aws/aws-sdk-go-v2/config v1.32.39
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.63.5
	go.uber.org/zap v1.28.0
	gopkg.in/guregu/null.v3 v3.5.0
)

require (
	github.com/aws/aws-sdk-go-v2/credentials v1.19.38 // indirect
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.20.62 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.18.39 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.39 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.39 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.6 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.4.40 // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodbstreams v1.36.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.18 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.12.16 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.39 // indirect
	github.com/aws/aws-sdk-go-v2/service/signin v1.5.8 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.33.8 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.38.8 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.45.8 // indirect
	github.com/aws/smithy-go v1.27.10 // indirect
	go.uber.org/multierr v1.10.0 // indirect
)

replace github.com/Kaaaaazuya/rss-commentator/go/shared v0.0.0 => ../../shared
