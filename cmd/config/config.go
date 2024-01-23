package config

import (
	"github.com/caarlos0/env/v8"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	PostgresDBAddr	string `env:"DATABASE_DSN,file" envDefault:"/c/Dev/EM"`
}

func ConfigSetup() Config {

	cfg := Config{}
	err := env.Parse(&cfg)
    if err != nil {
        log.Fatal(err)
    }

	return cfg
}