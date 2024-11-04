package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type key string

const KeyUUID key = key("uuid")

type Config struct {
	Service  Service
	Postgres Postgres
}

type Service struct {
	Port string `env:"SOCIETY_SERVICE_PORT"`
	Host string `env:"SOCIETY_SERVICE_HOST"`
}

type Postgres struct {
	User     string `env:"SOCIETY_SERVICE_POSTGRES_USER"`
	Password string `env:"SOCIETY_SERVICE_POSTGRES_PASSWORD"`
	Database string `env:"SOCIETY_SERVICE_POSTGRES_DB"`
	Host     string `env:"SOCIETY_SERVICE_POSTGRES_HOST"`
	Port     string `env:"SOCIETY_SERVICE_POSTGRES_PORT"`
}

func MustLoad() *Config {
	cfg := &Config{}
	err := cleanenv.ReadEnv(cfg)

	if err != nil {
		log.Fatalf("Can not read env variables: %s", err)
	}

	return cfg
}
