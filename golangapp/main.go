package main

import (
	"context"
	"net/http"
	"os"

	"github.com/angadthandi/golangmongoapp/golangapp/config"
	"github.com/angadthandi/golangmongoapp/golangapp/messages"

	log "github.com/sirupsen/logrus"

	// https://godoc.org/github.com/mongodb/mongo-go-driver/mongo
	"github.com/mongodb/mongo-go-driver/mongo"
)

var (
	// RabbitMQ messaging client
	MessagingClient messages.IMessagingClient

	// Messaging registry client
	MessagesRegistryClient messages.IMessagesRegistry
)

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

	dbUrl := "mongodb://" +
		config.MongoDBUsername + ":" +
		config.MongoDBPassword + "@" +
		config.MongoDBServiceName +
		config.MongoDBPort
	dbClient, err := mongo.Connect(context.Background(), dbUrl, nil)
	if err != nil {
		log.Fatalf("mongodb connection error : %v", err)
	}

	defer dbClient.Disconnect(context.Background())

	// connect to RabbitMQ
	MessagingClient = &messages.MessagingClient{}
	MessagingClient.Connect()

	defer MessagingClient.Close()

	// initialize message registry map
	MessagesRegistryClient = &messages.MessagesRegistryClient{}
	MessagesRegistryClient.InitCorrelationMap()

	// start receiver
	// listen to messages from RabbitMQ
	// sent by other micro services
	go MessagingClient.Receive(
		config.ExchangeName,
		config.ExchangeType,
		config.GoappRoutingKey,
		// events.HandleRefreshEvent,
		configureMessageRoutes,
		MessagesRegistryClient,
	)

	// configure route handlers
	configureRoutes(dbClient)

	log.Printf("Listening on Port: %v", config.ServerPort)
	// start http web server
	log.Fatal(http.ListenAndServe(config.ServerPort, nil))
}
