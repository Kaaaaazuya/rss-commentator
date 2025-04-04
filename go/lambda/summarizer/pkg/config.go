package pkg

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
)

type LLMConfig struct {
	Endpoint string
	Model    string
	APIKey   string
}

type Config struct {
	AWS aws.Config
	LLM *LLMConfig
}

func LoadConfig(ctx context.Context) (*Config, error) {
	awsCfg, err := awsConfig.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	llmCfg, err := loadLLMConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return &Config{
		AWS: awsCfg,
		LLM: llmCfg,
	}, nil
}

func loadLLMConfig() (*LLMConfig, error) {
	e := os.Getenv("LLM_ENDPOINT")
	if e == "" {
		return nil, fmt.Errorf("LLM_ENDPOINT is not set")
	}
	m := os.Getenv("LLM_MODEL")
	if m == "" {
		return nil, fmt.Errorf("LLM_MODEL is not set")
	}
	a := os.Getenv("LLM_API_KEY")
	if a == "" {
		return nil, fmt.Errorf("LLM_API_KEY is not set")
	}
	return &LLMConfig{
		Endpoint: e,
		Model:    m,
		APIKey:   a,
	}, nil
}
