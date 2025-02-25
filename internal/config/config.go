package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Service  Service
	Postgres Postgres
	Platform Platform
	Logger   Logger
}

type Service struct {
	Port string `env:"SOCIETY_SERVICE_PORT"`
	Host string `env:"SOCIETY_SERVICE_HOST"`
	Name string `env:"SOCIETY_SERVICE_NAME"`
}

type Postgres struct {
	User     string `env:"SOCIETY_SERVICE_POSTGRES_USER"`
	Password string `env:"SOCIETY_SERVICE_POSTGRES_PASSWORD"`
	Database string `env:"SOCIETY_SERVICE_POSTGRES_DB"`
	Host     string `env:"SOCIETY_SERVICE_POSTGRES_HOST"`
	Port     string `env:"SOCIETY_SERVICE_POSTGRES_PORT"`
}

type Platform struct {
	Env string `env:"ENV"`
}

type Logger struct {
	Host string `env:"LOGGER_SERVICE_HOST"`
	Port string `env:"LOGGER_SERVICE_PORT"`
}

func MustLoad() *Config {
	cfg := &Config{}
	err := cleanenv.ReadEnv(cfg)

	if err != nil {
		log.Fatalf("Can not read env variables: %s", err)
	}

	return cfg
}
