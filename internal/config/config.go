package config

import (
	"fmt"
	"os"

	"github.com/caarlos0/env/v6"
)

type API struct {
	Network        string `env:"API_NETWORK" envDefault:"tcp"`
	Address        string `env:"API_ADDRESS,required" example:"localhost:8082"`
	PublicHost     string `env:"API_PUBLIC_HOST,required"`
	MaxReceiveSize int    `env:"API_MAX_RECEIVE_SIZE" envDefault:"20"`
	MaxSendSize    int    `env:"API_MAX_SEND_SIZE" envDefault:"20"`
}

type S3 struct {
	AccessKey    string `env:"S3_ACCESS_KEY,required"`
	SecretKey    string `env:"S3_SECRET_KEY,required"`
	SessionToken string `env:"S3_SESSION_TOKEN" envDefault:""`
	Bucket       string `env:"S3_BUCKET" envDefault:"media"`
	Address      string `env:"S3_ADDRESS,required"`
	PolicyFile   string `env:"S3_POLICY_FILE,required"`
	Secure       bool   `env:"S3_SECURE" envDefault:"false"`
	Policy       string
}

type Metrics struct {
	Address string `env:"METRICS_ADDRESS,required" example:":8080"`
}

type Config struct {
	API     API
	S3      S3
	Metrics Metrics
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
