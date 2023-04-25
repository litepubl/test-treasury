package app

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Name string `yaml:"name" env:"NAME" envDefault:"test-player"`
	Env  string `yaml:"env" env:"ENV" envDefault:"development"`
	PG   struct {
		PoolMax int    `env-required:"true" yaml:"pool_max" env:"PG_POOL_MAX"`
		URL     string `env-required:"true"                 env:"PG_URL"`
	} `yaml:"postgres" env:"postgres"`
		HTTP struct {
	Port string    `yaml:"port" env:"PORT" envDefault:"8080"`
} `yaml:"http" env:"http"`

	Log  struct {
		Path      string `yaml:"path" env:"LOG_PATH" envDefault:"/app/log/"`
		ErrorFile string `yaml:"error_file" env:"LOG_ERROR_FILE" envDefault:"error.log"`
		DebugFile string `yaml:"debug_file" env:"LOG_DEBUG_FILE" envDefault:"debug.log"`
	} `yaml:"log" env:"LOG"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, err
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
