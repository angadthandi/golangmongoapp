package main

import (
	"encoding/json"

	"github.com/angadthandi/golangmongoapp/golangapp/config"
	"github.com/angadthandi/golangmongoapp/golangapp/jsondefinitions"
	"github.com/angadthandi/golangmongoapp/golangapp/messages"
	"github.com/mongodb/mongo-go-driver/mongo"
	log "github.com/sirupsen/logrus"
)

// IWriteReplyTo mocks Hub based on implemented funcs
type IWriteReplyTo interface {
	SendMsgToAllClients(jsonMsg json.RawMessage)
}

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
		log.Errorf(`golangapp configureMessageRoutes:
			unable to unmarshal json: %v`, err)

		SendErrorResponse(
			MessagesRegistryClient,
			replyToRoutingKey,
			receivedCorrelationId,
			isReplyMessage,
		)
		return
	}

	var (
		resp interface{}

		// set this based on incoming Type
		// in the switch statement
		//
		// for some cases, we might need to
		// send message to all clients
		// for other cases, send to specific client
		sendMsgToAllClients bool
	)

	if isResponseToExistingMessage {
		// Handle Response to Sent Message Switch Cases

		log.Debugf(`golangapp configureMessageRoutes:
		msg.Type %v`, msg.Type)
		log.Debugf(`golangapp configureMessageRoutes:
		msg.Message %v`, msg.Message)
		switch msg.Type {
		case "GetProducts":
			resp = msg.Message
			sendMsgToAllClients = false

		default:
			resp = msg.Message
			sendMsgToAllClients = true
		}
	} else {
		// Handle New Incoming Message Switch Cases

		// As this app is the API Gateway
		// we might not need any cases here
	}

	// write response messages to
	// WS
	log.Debugf("golangapp configureMessageRoutes resp: %v", resp)

	if writeReplyTo != nil {
		i, ok := writeReplyTo.(IWriteReplyTo)
		if !ok {
			log.Errorf(`golangapp configureMessageRoutes invalid
			writeReplyTo: %v, of Type: %T`, resp, resp)
			return
		}

		iResp, ok := resp.(json.RawMessage)
		if !ok {
			log.Errorf(`golangapp configureMessageRoutes invalid
			iResp: %v, of Type: %T`, iResp, iResp)
			return
		}

		if sendMsgToAllClients {
			log.Debug(`golangapp configureMessageRoutes
			 SendMsgToAllClients`)
			// send to all connected clients
			i.SendMsgToAllClients(iResp)
		} else {
			log.Debug(`golangapp configureMessageRoutes
			 SendMsgToClientWithCorrelationId`)
			// send to client with correlationId
			// in client's clientCorrelationIds map

			if hubCh == nil {
				log.Debug(`golangapp configureMessageRoutes
					 nil hubCh`)
				return
			}

			var chMsg jsondefinitions.MicroServiceResponseMsgForHub

			chMsg.CorrelationId = receivedCorrelationId
			chMsg.ReceivedJsonMsg = iResp
			log.Debugf("configureMessageRoutes send on hubCh chMsg: %v",
				chMsg)

			b, err := json.Marshal(chMsg)
			if err != nil {
				log.Errorf("configureMessageRoutes unable to marshal: %v",
					err)
				return
			}

			hubCh <- b
		}
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
	msg.Type = responseType
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
