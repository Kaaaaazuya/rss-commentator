package pkg

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
)

type Config struct {
	AWS aws.Config
}

func LoadConfig(ctx context.Context) (*Config, error) {
	awsCfg, err := awsConfig.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return &Config{
		AWS: awsCfg,
	}, nil
}
