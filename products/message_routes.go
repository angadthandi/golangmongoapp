package main

import (
	"encoding/json"

	"github.com/angadthandi/golangmongoapp/products/dummySend"
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

	if !isResponseToExistingMessage {
		dummySend.DummySendToGoapp(
			MessagingClient,
			MessagesRegistryClient,
			replyToRoutingKey,
			receivedCorrelationId,
			isReplyMessage,
		)
	}
}
