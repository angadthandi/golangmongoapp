package main

import (
	"os"

	"github.com/angadthandi/golangmongoapp/products/config"
	"github.com/angadthandi/golangmongoapp/products/events"
	"github.com/angadthandi/golangmongoapp/products/messages"

	log "github.com/sirupsen/logrus"
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
	log.Infof("Started main: %v", "products")

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
	MessagingClient.Receive(
		config.ExchangeName,
		config.ExchangeType,
		config.ProductsRoutingKey,
		events.HandleRefreshEvent,
		MessagesRegistryClient,
	)
}
