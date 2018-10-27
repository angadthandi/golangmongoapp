package main

import (
	"net/http"

	"github.com/angadthandi/golangmongoapp/golangapp/config"
	"github.com/angadthandi/golangmongoapp/golangapp/messages"

	mgo "gopkg.in/mgo.v2"
)

// RabbitMQ messaging client
var MessagingClient messages.IMessagingClient

func main() {
	// initialize logging
	initLogger()
	log.Info("Starting goapp main...")

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
		[]string{config.ProductsRoutingKey},
	)

	// configure route handlers
	configureRoutes(dbSession)

	log.Printf("Listening on Port: %v", config.ServerPort)
	// start http web server
	log.Fatal(http.ListenAndServe(config.ServerPort, nil))
}
