package main

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func loggerInit() {
	if config.Environment == "dev" {
		log.SetOutput(os.Stdout)
		log.SetFormatter(&log.TextFormatter{})
		log.SetLevel(log.DebugLevel)
	} else {
		 file, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY, 0666)
		 if err == nil {
		  	log.SetOutput(file)
		 } else {
		  	log.Info("Failed to log to file, using default stderr")
		 }
		log.SetFormatter(&log.JSONFormatter{})
		log.SetLevel(log.ErrorLevel)
	}
}
