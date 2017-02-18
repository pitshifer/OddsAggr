package main

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func loggerInit(cfg Config) {
	log.SetOutput(os.Stdout)
	if cfg.Environment == "dev" {
		log.SetFormatter(&log.TextFormatter{})
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetFormatter(&log.JSONFormatter{})
		log.SetLevel(log.ErrorLevel)
	}
}
