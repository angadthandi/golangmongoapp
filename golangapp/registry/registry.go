package registry

import (
	"encoding/json"

	"github.com/angadthandi/golangmongoapp/golangapp/api/ws"
	"github.com/angadthandi/golangmongoapp/golangapp/messages"
	// "github.com/mongodb/mongo-go-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo"

	// This pakcage cannot inlcude goappsocket package
	// otherwise it results in a circular dependency
	// "github.com/angadthandi/golangmongoapp/golangapp/goappsocket"
	log "github.com/sirupsen/logrus"
)

func ClientRegistry(
	c ws.Clienter,
	dbClient *mongo.Client,
	MessagingClient messages.IMessagingClient,
	MessagesRegistryClient messages.IMessagesRegistry,
	jsonMsg json.RawMessage,
	ChClientCorrelationIds chan<- []byte,
) {
	log.Debugf("registry ClientRegistry")
	ws.ClientAPI(
		c,
		dbClient,
		MessagingClient,
		MessagesRegistryClient,
		jsonMsg,
		ChClientCorrelationIds,
	)
}

func HubRegistry(
	h ws.Huber,
	dbClient *mongo.Client,
	MessagingClient messages.IMessagingClient,
	MessagesRegistryClient messages.IMessagesRegistry,
	jsonMsg json.RawMessage,
) {
	log.Debugf("registry HubRegistry")
	ws.HubAPI(
		h,
		dbClient,
		MessagingClient,
		MessagesRegistryClient,
		jsonMsg,
	)
}
