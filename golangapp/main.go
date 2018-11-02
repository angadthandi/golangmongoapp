package main

import (
	"net/http"
	"os"

	"github.com/angadthandi/golangmongoapp/golangapp/config"
	"github.com/angadthandi/golangmongoapp/golangapp/events"
	"github.com/angadthandi/golangmongoapp/golangapp/messages"

	log "github.com/sirupsen/logrus"
	mgo "gopkg.in/mgo.v2"
)

// RabbitMQ messaging client
var MessagingClient messages.IMessagingClient

// initialize logger
func init() {
	// Log as JSON instead of the default ASCII formatter.
	// log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
}

func main() {
	log.Infof("Started main: %v", "goapp")

	// connect to database
	dbSession, err := mgo.Dial(
		"mongodb://" +
			config.MongoDBUsername + ":" +
			config.MongoDBPassword + "@" +
			config.MongoDBServiceName +
			config.MongoDBPort)
	if err != nil {
		log.Fatalf("mongodb connection error : %v", err)
	}

	defer dbSession.Close()

	dbSession.SetMode(mgo.Monotonic, true)

	// connect to RabbitMQ
	MessagingClient = &messages.MessagingClient{}
	MessagingClient.Connect()

	defer MessagingClient.Close()

	// start receiver
	// listen to messages from RabbitMQ
	// sent by other micro services
	go MessagingClient.Receive(
		config.ExchangeName,
		config.ExchangeType,
		config.GoappRoutingKey,
		events.HandleRefreshEvent,
	)

	// configure route handlers
	configureRoutes(dbSession)

	log.Printf("Listening on Port: %v", config.ServerPort)
	// start http web server
	log.Fatal(http.ListenAndServe(config.ServerPort, nil))
}
