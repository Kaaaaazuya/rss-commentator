package config

import "github.com/aws/aws-sdk-go-v2/aws"

type EnvConfig struct {
	App *App
	AWS aws.Config
}

// LoadEnvConfig loads the environment variables and returns the configuration
func LoadEnvConfig() (*EnvConfig, error) {
	app, err := loadAppConfig()
	if err != nil {
		return nil, err
	}

	awsConf, err := loadAWSConfig()
	if err != nil {
		return nil, err
	}

	c := &EnvConfig{
		App: app,
		AWS: awsConf,
	}

	return c, nil
}
