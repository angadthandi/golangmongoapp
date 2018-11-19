package main

import (
	"encoding/json"

	"github.com/mongodb/mongo-go-driver/mongo"

	"github.com/angadthandi/golangmongoapp/products/config"
	"github.com/angadthandi/golangmongoapp/products/controllers"
	"github.com/angadthandi/golangmongoapp/products/jsondefinitions"
	"github.com/angadthandi/golangmongoapp/products/messages"
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
	writeReplyTo interface{}, // write reply to http/ws
	hubCh chan []byte,
) {
	log.Debugf(`products configureMessageRoutes:
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
	log.Debugf(`products configureMessageRoutes:
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

		var resp interface{}

		log.Debugf(`products configureMessageRoutes:
		msg.Type %v`, msg.Type)
		switch msg.Type {
		case "GetProducts":
			// Add Product for Testing GetProducts
			// controllers.CreateProduct(
			// 	dbRef,
			// 	"Test Product 1",
			// 	"Test Product Code 1",
			// )
			resp = controllers.GetProducts(dbRef)

		default:
			resp = "Invalid Message Type"
		}

		SendSuccessResponse(
			MessagesRegistryClient,
			replyToRoutingKey,
			receivedCorrelationId,
			isReplyMessage,
			resp,
			msg.Type,
		)
		// dummySend.DummySendToGoapp(
		// 	MessagingClient,
		// 	MessagesRegistryClient,
		// 	replyToRoutingKey,
		// 	receivedCorrelationId,
		// 	isReplyMessage,
		// )
	}
}

func SendSuccessResponse(
	MessagesRegistryClient messages.IMessagesRegistry,
	sendToRoutingKey string,
	receivedCorrelationId string,
	isReplyMessage bool,
	responseMessage interface{},
	responseType string,
) {
	var msg jsondefinitions.GenericMessageSend
	msg.Type = responseType //"Success"
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
		log.Errorf("Products: SendErrorResponse: unable to marshal: %v",
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

	log.Debug("Products SendErrorResponse Sent!")
}
