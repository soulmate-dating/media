package config

import (
	"fmt"
	"os"

	"github.com/caarlos0/env/v6"
)

type API struct {
	Network string `env:"API_NETWORK" envDefault:"tcp"`
	Address string `env:"API_ADDRESS" envDefault:"8082"`
}

type S3 struct {
	AccessKey    string `env:"S3_ACCESS_KEY"`
	SecretKey    string `env:"S3_SECRET_KEY"`
	SessionToken string `env:"S3_SESSION_TOKEN" envDefault:""`
	Bucket       string `env:"S3_BUCKET" envDefault:"media"`
	Address      string `env:"S3_ADDRESS"`
	PolicyFile   string `env:"S3_POLICY_FILE"`
	Secure       bool   `env:"S3_SECURE" envDefault:"false"`
	Policy       string
}

type Config struct {
	API API
	S3  S3
}

func Load() (Config, error) {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		return Config{}, fmt.Errorf("parsing config: %w", err)
	}
	if cfg.S3.PolicyFile != "" {
		policy, err := os.ReadFile(cfg.S3.PolicyFile)
		if err != nil {
			return Config{}, fmt.Errorf("reading policy file: %w", err)
		}
		cfg.S3.Policy = string(policy)
	}

	return cfg, nil
}
