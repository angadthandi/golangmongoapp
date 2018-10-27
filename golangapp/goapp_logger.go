package main

import (
	"os"

	logrus "github.com/sirupsen/logrus"
)

var log *logrus.Logger

func initLogger() {
	log = logrus.New()

	// file, err := os.OpenFile("goapp_logs.log", os.O_CREATE|os.O_WRONLY, 0666)
	// if err == nil {
	// 	log.Out = file
	// } else {
	// 	log.Info("Failed to log to file, using default stderr")
	// }

	log.SetOutput(os.Stdout)
	// log.SetFormatter(&logrus.JSONFormatter{})
	// log.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	log.SetLevel(logrus.DebugLevel)

	// logrus types
	// log.Trace("Something very low level.")
	// log.Debug("Useful debugging information.")
	// log.Info("Something noteworthy happened!")
	// log.Warn("You should probably take a look at this.")
	// log.Error("Something failed but I'm not quitting.")
}
