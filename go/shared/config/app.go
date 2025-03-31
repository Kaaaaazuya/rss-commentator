package config

import "github.com/kelseyhightower/envconfig"

type App struct {
	Host string `required:"true"`
	Env  string `required:"true"`
}

func loadAppConfig() (*App, error) {
	var cfg App
	err := envconfig.Process("app", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
