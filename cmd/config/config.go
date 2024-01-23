package config

import (
	log "github.com/sirupsen/logrus"
	"os"
)

type Config struct {
	PostgresDBAddr	string 
}

func ConfigSetup() Config {
	cfg := Config{}
	var exists bool
	cfg.PostgresDBAddr, exists = os.LookupEnv("DATABASE_DSN")

	if exists!= true {
		log.Fatal("нет пути к базе данных!")
	}

	return cfg
}