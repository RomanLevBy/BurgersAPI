package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string `yaml:"env" env:"ENV" env-required:"true"`
	Postgres   `yaml:"postgres"`
	HTTPServer `yaml:"http_server"`
}

type Postgres struct {
	Host     string `env:"POSTGRES_CONTAINER_NAME" env-default:"burgers-api-postgres-db"`
	User     string `env:"POSTGRES_USER" env-required:"true"`
	Password string `env:"POSTGRES_PASSWORD" env-required:"true"`
	DBName   string `env:"POSTGRES_DB" env-default:"burgers_api"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func Load() (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		return nil, errors.New("CONFIG_PATH is not set")
	}

	//check if file exist
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, errors.New(fmt.Sprintf("config file %s does not exist", configPath))
	}

	var conf Config

	if err := cleanenv.ReadConfig(configPath, &conf); err != nil {
		return nil, errors.New(fmt.Sprintf("Cannot read config file: %s", err))
	}

	return &conf, nil
}
