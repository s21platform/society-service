package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type Config struct {
	Service Service
}

type Service struct {
	Port string `env:"SOCIETY_SERVICE_PORT"`
	Host string `env:"SOCIETY_SERVICE_HOST"`
}

type ReadEnvBD struct {
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
