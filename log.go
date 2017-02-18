package main

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func loggerInit() {
	log.SetOutput(os.Stdout)
	if config.Environment == "dev" {
		log.SetFormatter(&log.TextFormatter{})
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetFormatter(&log.JSONFormatter{})
		log.SetLevel(log.ErrorLevel)
	}
}
