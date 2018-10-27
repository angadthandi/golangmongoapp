package main

import (
	"os"

	logrus "github.com/sirupsen/logrus"
)

var log *logrus.Logger

func initLogger() {
	log = logrus.New()
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.DebugLevel)
}
