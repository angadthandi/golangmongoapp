package main

import (
	"context"
	"os"

	"github.com/angadthandi/golangmongoapp/products/config"
	"github.com/angadthandi/golangmongoapp/products/messages"
	"github.com/mongodb/mongo-go-driver/mongo"

	log "github.com/sirupsen/logrus"
)

const ProductsDBName = "productsdb"

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

	// connect to products mongodb
	dbUrl := "mongodb://" +
		config.MongoDBUsername + ":" +
		config.MongoDBPassword + "@" +
		config.MongoDBServiceName +
		config.MongoDBPort
	dbClient, err := mongo.Connect(context.Background(), dbUrl, nil)
	if err != nil {
		log.Fatalf("mongodb connection error : %v", err)
	}

	log.Debug("Connected to mongodb products database")
	defer dbClient.Disconnect(context.Background())

	dbRef := dbClient.Database(ProductsDBName)
	log.Debug("Initialized mongodb products database")

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
		configureMessageRoutes,
		MessagesRegistryClient,
		dbRef,
	)
}
