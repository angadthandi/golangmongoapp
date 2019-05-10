package ws

import (
	"encoding/json"

	"github.com/angadthandi/golangmongoapp/golangapp/config"
	"github.com/angadthandi/golangmongoapp/golangapp/jsondefinitions"
	"github.com/angadthandi/golangmongoapp/golangapp/messages"
	// "github.com/mongodb/mongo-go-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo"
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
	GetClientUUID() string
}

type Huber interface {
	Run()
	SendMsgToAllClients(jsonMsg json.RawMessage)
	SendMsgToClientWithCorrelationId(
		jsonMsg json.RawMessage,
		correlationId string,
	)
}

// handler for ws/API
func ClientAPI(
	c Clienter,
	dbClient *mongo.Client,
	MessagingClient messages.IMessagingClient,
	MessagesRegistryClient messages.IMessagesRegistry,
	jsonMsg json.RawMessage,
	ChClientCorrelationIds chan<- []byte,
) {
	log.Debug("ws ClientAPI")

	respMsg, correlationId := API(
		dbClient,
		MessagingClient,
		MessagesRegistryClient,
		jsonMsg,
	)
	log.Debugf("ws ClientAPI respMsg: %v\n", respMsg)
	log.Debugf("ws ClientAPI correlationId: %v\n", correlationId)

	if correlationId != "" {
		var clientUUIDCorrId jsondefinitions.ClientUUIDCorrelationID

		// get clientUUID
		clientUUIDCorrId.ClientUUID = c.GetClientUUID()
		clientUUIDCorrId.ClientCorrelationId = correlationId

		b, err := json.Marshal(clientUUIDCorrId)
		if err != nil {
			log.Debugf("ws ClientAPI Unable to marshal: %v",
				err)
			return
		}

		ChClientCorrelationIds <- b
	} else {
		// send message to all clients on hub
		c.SendMessageOnHub(respMsg)
	}
}

// handler for ws/API
func HubAPI(
	h Huber,
	dbClient *mongo.Client,
	MessagingClient messages.IMessagingClient,
	MessagesRegistryClient messages.IMessagesRegistry,
	jsonMsg json.RawMessage,
) {
	log.Debugf("ws HubAPI")

	respMsg, correlationId := API(
		dbClient,
		MessagingClient,
		MessagesRegistryClient,
		jsonMsg,
	)
	log.Debugf("ws HubAPI correlationId: %v\n", correlationId)

	h.SendMsgToAllClients(respMsg)
}

// handler for ws/API
func API(
	dbClient *mongo.Client,
	MessagingClient messages.IMessagingClient,
	MessagesRegistryClient messages.IMessagesRegistry,
	jsonMsg json.RawMessage,
) (json.RawMessage, string) {
	var msg jsondefinitions.GenericAPIRecieve
	err := json.Unmarshal(jsonMsg, &msg)
	if err != nil {
		log.Errorf("ws/API JSON Unmarshal error: %v", err)
		return nil, ""
	}

	var correlationId string
	log.Debugf("ws/API msg.Api: %v", msg.Api)
	switch msg.Api {
	case "GetProducts":
		// get data from products service
		correlationId = sendToMessageQueue(
			MessagingClient,
			MessagesRegistryClient,
			config.ProductsRoutingKey,
			msg.Api,
			msg.Message,
		)
	case "CreateProduct":
		correlationId = sendToMessageQueue(
			MessagingClient,
			MessagesRegistryClient,
			config.ProductsRoutingKey,
			msg.Api,
			msg.Message,
		)

	default:
		correlationId = ""
	}

	var resp jsondefinitions.GenericAPIResponse

	resp.Api = msg.Api
	resp.Message = msg.Message

	b, err := json.Marshal(resp)
	if err != nil {
		log.Errorf("ws/API JSON Marshal error: %v", err)
		return nil, ""
	}

	log.Debugf("ws/API JSON Response: %v", string(b))
	return b, correlationId
}

func sendToMessageQueue(
	MessagingClient messages.IMessagingClient,
	MessagesRegistryClient messages.IMessagesRegistry,
	publishRoutingKey string,
	msgType string,
	msgData interface{},
) string {
	var m jsondefinitions.GenericMessageSend
	m.Type = msgType
	m.Message = msgData

	b, err := json.Marshal(m)
	if err != nil {
		log.Errorf("sendToMessageQueue unable to marshal: %v", err)
	}

	correlationId := MessagingClient.Send(
		config.ExchangeName,
		config.ExchangeType,
		publishRoutingKey, //config.ProductsRoutingKey,
		config.GoappRoutingKey,
		b,
		MessagesRegistryClient,
		"",
		false,
	)

	return correlationId
}
