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

func MustLoad() *Config {
	cfg := &Config{}
	err := cleanenv.ReadEnv(cfg)

	if err != nil {
		log.Fatalf("Can not read env variables: %s", err)
	}

	return cfg
}
