package main

import (
	"encoding/json"

	"github.com/angadthandi/golangmongoapp/golangapp/config"
	"github.com/angadthandi/golangmongoapp/golangapp/jsondefinitions"
	"github.com/angadthandi/golangmongoapp/golangapp/messages"
	"github.com/mongodb/mongo-go-driver/mongo"
	log "github.com/sirupsen/logrus"
)

// configure RabbitMQ Message Routes
func configureMessageRoutes(
	MessagingClient *messages.MessagingClient,
	MessagesRegistryClient messages.IMessagesRegistry,
	jsonMsg json.RawMessage,
	replyToRoutingKey string,
	receivedCorrelationId string,
	isResponseToExistingMessage bool,
	dbRef *mongo.Database,
) {
	log.Debugf(`golangapp configureMessageRoutes:
	jsonMsg: %v, replyToRoutingKey: %v,
	receivedCorrelationId: %v, isResponseToExistingMessage: %v`,
		jsonMsg, replyToRoutingKey,
		receivedCorrelationId, isResponseToExistingMessage)

	// isReplyMessage should be the OPPOSITE of isResponseToExistingMessage.
	// Variable assigned for verbosity
	//
	// message_events: handleRefreshEvent
	// sets isResponseToExistingMessage
	// based on if a correlationId was found in our current app
	isReplyMessage := !isResponseToExistingMessage
	log.Debugf(`golangapp configureMessageRoutes:
		isReplyMessage %v`, isReplyMessage)

	var msg jsondefinitions.GenericMessageRecieve

	err := json.Unmarshal(jsonMsg, &msg)
	if err != nil {
		log.Errorf(`products configureMessageRoutes:
			unable to unmarshal json: %v`, err)

		SendErrorResponse(
			MessagesRegistryClient,
			replyToRoutingKey,
			receivedCorrelationId,
			isReplyMessage,
		)
		return
	}

	if isResponseToExistingMessage {
		// Handle Response to Sent Message Switch Cases
	} else {
		// Handle New Incoming Message Switch Cases
	}
}

func SendSuccessResponse(
	MessagesRegistryClient messages.IMessagesRegistry,
	sendToRoutingKey string,
	receivedCorrelationId string,
	isReplyMessage bool,
	responseMessage interface{},
) {
	var msg jsondefinitions.GenericMessageSend
	msg.Type = "Success"
	msg.Message = responseMessage

	b, err := json.Marshal(msg)
	if err != nil {
		// This should not happen.
		// If we are here
		// then response wont be sent to a waiting microservice.
		// That microservice will have a hanging correlationId
		// in its map.
		log.Errorf("Products: SendSuccessResponse: unable to marshal: %v",
			err)
		return
	}

	MessagingClient.Send(
		config.ExchangeName,
		config.ExchangeType,
		sendToRoutingKey,
		config.ProductsRoutingKey,
		b,
		MessagesRegistryClient,
		receivedCorrelationId,
		isReplyMessage,
	)

	log.Debug("Products SendSuccessResponse Sent!")
}

func SendErrorResponse(
	MessagesRegistryClient messages.IMessagesRegistry,
	sendToRoutingKey string,
	receivedCorrelationId string,
	isReplyMessage bool,
) {
	var msg jsondefinitions.GenericMessageSend
	msg.Type = "Error"
	msg.Message = jsondefinitions.GenericErrorMessageSend{
		Errormessage: "Invalid JSON message",
	}

	b, err := json.Marshal(msg)
	if err != nil {
		// This should not happen.
		// If we are here
		// then response wont be sent to a waiting microservice.
		// That microservice will have a hanging correlationId
		// in its map.
		log.Errorf("Golangapp: SendErrorResponse: unable to marshal: %v",
			err)
		return
	}

	MessagingClient.Send(
		config.ExchangeName,
		config.ExchangeType,
		sendToRoutingKey,
		config.ProductsRoutingKey,
		b,
		MessagesRegistryClient,
		receivedCorrelationId,
		isReplyMessage,
	)

	log.Debug("Golangapp SendErrorResponse Sent!")
}
