package config

import (
	"os"

	errorsUtils "github.com/Egor213/notifyBot/pkg/errors"

	"github.com/ilyakaznacheev/cleanenv"
	log "github.com/sirupsen/logrus"

	"github.com/joho/godotenv"
)

type (
	Config struct {
		App `yaml:"app"`
		Log `yaml:"log"`
		PG  `yaml:"postgres"`
	}

	App struct {
		Name    string `yaml:"name" env-required:"true"`
		Version string `yaml:"version" env-required:"true"`
	}

	Log struct {
		Level string `yaml:"level" env:"LOG_LEVEL" env-default:"info"`
	}

	PG struct {
		MaxPoolSize int    `env-required:"true" env:"MAX_POOL_SIZE" yaml:"max_pool_size"`
		URL         string `env-required:"true" env:"PG_URL"`
	}
)

const ENV_PATH = "infra/.env"

func init() {
	if err := godotenv.Load(ENV_PATH); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func New() (*Config, error) {
	cfg := &Config{}

	pathToConfig, ok := os.LookupEnv("APP_CONFIG_PATH")
	if !ok || pathToConfig == "" {
		log.WithField("env_var", "APP_CONFIG_PATH").
			Info("Config path is not set, using default")
		pathToConfig = "infra/config.yaml"
	}

	if err := cleanenv.ReadConfig(pathToConfig, cfg); err != nil {
		return nil, errorsUtils.WrapPathErr(err)
	}

	if err := cleanenv.UpdateEnv(cfg); err != nil {
		return nil, errorsUtils.WrapPathErr(err)
	}

	return cfg, nil
}
