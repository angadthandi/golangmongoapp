package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/angadthandi/golangmongoapp/golangapp/config"
	"github.com/angadthandi/golangmongoapp/golangapp/goappsocket"
	"github.com/angadthandi/golangmongoapp/golangapp/messages"

	log "github.com/sirupsen/logrus"

	// https://godoc.org/github.com/mongodb/mongo-go-driver/mongo
	// "github.com/mongodb/mongo-go-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	// dbClient, err := mongo.Connect(context.Background(), dbUrl, nil)
	// dbClient, err := mongo.NewClient(options.Client().ApplyURI(dbUrl))
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	dbClient, err := mongo.Connect(ctx, options.Client().ApplyURI(dbUrl))
	if err != nil {
		log.Fatalf("mongodb connection error : %v", err)
	}

	log.Debug("Connected to mongodb golangapp database")
	defer dbClient.Disconnect(context.Background())

	dbRef := dbClient.Database("GolangappDB")
	log.Debug("Initialized mongodb golangapp database")

	// connect to RabbitMQ
	MessagingClient = &messages.MessagingClient{}
	MessagingClient.Connect()

	defer MessagingClient.Close()

	// initialize message registry map
	MessagesRegistryClient = &messages.MessagesRegistryClient{}
	MessagesRegistryClient.InitCorrelationMap()

	// start hub
	// for creating websocket conns
	//
	// pass this hub to the message Receiver
	// so we can send responses from other services
	// to the websocket
	hub := goappsocket.NewHub()
	go hub.Run()

	// configure route handlers
	configureRoutes(
		dbClient,
		hub,
		MessagingClient,
		MessagesRegistryClient,
	)

	// start receiver
	// listen to messages from RabbitMQ
	// sent by other micro services
	go MessagingClient.Receive(
		config.ExchangeName,
		config.ExchangeType,
		config.GoappRoutingKey,
		configureMessageRoutes,
		MessagesRegistryClient,
		dbRef,
		hub,
		hub.ChSendMsgToClientWithCorrelationId,
	)

	log.Printf("Listening on Port: %v", config.ServerPort)
	// start http web server
	log.Fatal(http.ListenAndServe(config.ServerPort, nil))
}
