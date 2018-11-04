package events

import (
	"encoding/json"

	"github.com/angadthandi/golangmongoapp/golangapp/messages"
	"github.com/angadthandi/golangmongoapp/golangapp/messagesRegistry"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

// {"type":"RefreshRemoteApplicationEvent",
// "timestamp":1494514362123,
// "originService":"config-server:docker:8888",
// "destinationService":"xxxaccoun:**",
// "id":"53e61c71-cbae-4b6d-84bb-d0dcc0aeb4dc"}
type UpdateToken struct {
	Type               string `json:"type"`
	Timestamp          int    `json:"timestamp"`
	OriginService      string `json:"originService"`
	DestinationService string `json:"destinationService"`
	Id                 string `json:"id"`
}

func HandleRefreshEvent(
	d amqp.Delivery,
	MessagingClient *messages.MessagingClient,
	MessagesRegistryClient messagesRegistry.IMessagesRegistry,
) {
	body := d.Body
	consumerTag := d.ConsumerTag
	correlationId := d.CorrelationId
	updateToken := &UpdateToken{}
	err := json.Unmarshal(body, updateToken)
	if err != nil {
		log.Printf("Problem parsing UpdateToken: %v", err.Error())
	} else {
		log.Debugf("HandleRefreshEvent: Received CorrelationId: %s", correlationId)
		log.Debugf("HandleRefreshEvent: Received ConsumerTag: %s", consumerTag)
		log.Debugf("HandleRefreshEvent: Received message: %s", body)

		// check if correlationId exists
		received, ok := MessagesRegistryClient.GetCorrelationData(correlationId)
		if ok {
			// existing event response by Outside App
			// handle response

			log.Debugf("HandleRefreshEvent: SentToAppName: %s, SentToAppEvent: %s",
				received.SentToAppName, received.SentToAppEvent)

			log.Debugf("HandleRefreshEvent: DeleteCorrelationMapData correlationId: %s",
				correlationId)
			MessagesRegistryClient.DeleteCorrelationMapData(correlationId)

			HandleResponseToExistingMessage(
				d.Body,
				MessagingClient,
				MessagesRegistryClient,
			)
		} else {
			log.Debug("HandleRefreshEvent: correlationId not found!")
			log.Debugf("HandleRefreshEvent: Received ReplyTo: %s", d.ReplyTo)

			// new event sent by Outside App
			// handle and respond back to Outside App
			HandleNewMessage(
				d.Body,
				MessagingClient,
				MessagesRegistryClient,
				d.ReplyTo,
				correlationId,
			)
		}

		// if strings.Contains(updateToken.DestinationService, consumerTag) {
		// 	log.Println("Consumertag is same as application name.")

		// 	// Consumertag is same as application name.

		// 	// https://github.com/callistaenterprise/goblog/blob/P9/common/config/loader.go
		// 	// LoadConfigurationFromBranch(
		// 	// 	viper.GetString("configServerUrl"),
		// 	// 	consumerTag,
		// 	// 	viper.GetString("profile"),
		// 	// 	viper.GetString("configBranch"))
		// }
	}
}

func HandleResponseToExistingMessage(
	jsonMsg json.RawMessage,
	MessagingClient *messages.MessagingClient,
	MessagesRegistryClient messagesRegistry.IMessagesRegistry,
) {
	log.Debugf("HandleResponseToExistingMessage")
	handleMessage(jsonMsg)
}

func HandleNewMessage(
	jsonMsg json.RawMessage,
	MessagingClient *messages.MessagingClient,
	MessagesRegistryClient messagesRegistry.IMessagesRegistry,
	replyToRoutingKey string,
	receivedCorrelationId string,
) {
	log.Debugf("HandleNewMessage")
	handleMessage(jsonMsg)
}

func handleMessage(
	jsonMsg json.RawMessage,
) {
	log.Debugf("handleMessage: Received JSON: %s",
		jsonMsg)
}
