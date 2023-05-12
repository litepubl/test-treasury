package app

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/litepubl/test-treasury/pkg/logger"
	"github.com/litepubl/test-treasury/pkg/postgres"
	"os"
)

const (
	DefaultConfigFilename = "./config/config.yml"
	AppConfigFileName     = "APP_CONFIG_FILENAME"
)

type Config struct {
	Name string          `yaml:"name" env:"NAME" envDefault:"test-player"`
	Env  string          `yaml:"env" env:"ENV" envDefault:"development"`
	PG   postgres.Config `yaml:"postgres" env:"postgres"`
	HTTP struct {
		Port string `yaml:"port" env:"PORT" envDefault:"8080"`
	} `yaml:"http" env:"http"`

	Log logger.Config `yaml:"log" env:"LOG"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{
		PG:  postgres.NewConfig(),
		Log: logger.Config{},
	}

	configFilename := DefaultConfigFilename
	env := os.Getenv(AppConfigFileName)
	if env != "" {
		configFilename = env
	}

	err := cleanenv.ReadConfig(configFilename, cfg)
	if err != nil {
		return nil, err
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
