package ws

import (
	"encoding/json"

	"github.com/angadthandi/golangmongoapp/golangapp/jsondefinitions"
	"github.com/angadthandi/golangmongoapp/golangapp/messages"
	"github.com/mongodb/mongo-go-driver/mongo"
	log "github.com/sirupsen/logrus"
)

// https://github.com/mantishK/dep
type Clienter interface {
	ReadPump(
		dbClient *mongo.Client,
		MessagingClient messages.IMessagingClient,
		MessagesRegistryClient messages.IMessagesRegistry,
	)
	WritePump()
	SendMessageOnHub(jsonMsg json.RawMessage)
}

type Huber interface {
	Run()
	SendMsgToAllClients(jsonMsg json.RawMessage)
}

// handler for ws/API
func ClientAPI(
	c Clienter,
	dbClient *mongo.Client,
	MessagingClient messages.IMessagingClient,
	MessagesRegistryClient messages.IMessagesRegistry,
	jsonMsg json.RawMessage,
) {
	log.Debugf("ws ClientAPI client Type: %T\n", c)
	log.Debugf("ws ClientAPI client: %v\n", c)

	respMsg := API(
		dbClient,
		MessagingClient,
		MessagesRegistryClient,
		jsonMsg,
	)

	c.SendMessageOnHub(respMsg)
	// log.Debugf("ws RegistryClient client Hub: %v", c.Hub)
}

// handler for ws/API
func HubAPI(
	h Huber,
	dbClient *mongo.Client,
	MessagingClient messages.IMessagingClient,
	MessagesRegistryClient messages.IMessagesRegistry,
	jsonMsg json.RawMessage,
) {
	log.Debugf("ws HubAPI hub Type: %T\n", h)
	log.Debugf("ws HubAPI hub: %v\n", h)

	respMsg := API(
		dbClient,
		MessagingClient,
		MessagesRegistryClient,
		jsonMsg,
	)

	h.SendMsgToAllClients(respMsg)
	// log.Debugf("ws RegistryClient client Hub: %v", c.Hub)
}

// handler for ws/API
func API(
	dbClient *mongo.Client,
	MessagingClient messages.IMessagingClient,
	MessagesRegistryClient messages.IMessagesRegistry,
	jsonMsg json.RawMessage,
) json.RawMessage {
	var msg jsondefinitions.GenericAPIRecieve
	err := json.Unmarshal(jsonMsg, &msg)
	if err != nil {
		log.Errorf("ws/API JSON Unmarshal error: %v", err)
		return nil
	}

	var resp jsondefinitions.GenericAPIResponse

	resp.Api = msg.Api
	// msg := "Test JSON Message!"
	resp.Message = msg.Message

	b, err := json.Marshal(resp)
	if err != nil {
		log.Errorf("ws/API JSON Marshal error: %v", err)
		return nil
	}

	log.Debugf("ws/API JSON Response: %v", string(b))
	return b
}
