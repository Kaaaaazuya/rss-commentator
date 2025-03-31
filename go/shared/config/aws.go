package config

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

func loadAWSConfig() (aws.Config, error) {
	awsConf, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return aws.Config{}, fmt.Errorf("load aws error. %w", err)
	}
	return awsConf, nil
}
