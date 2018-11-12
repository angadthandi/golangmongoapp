package messages

import (
	"encoding/json"

	"github.com/mongodb/mongo-go-driver/mongo"
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

func handleRefreshEvent(
	d amqp.Delivery,
	MessagingClient *MessagingClient,
	MessagesRegistryClient IMessagesRegistry,
	dbRef *mongo.Database,
	writeReplyTo interface{}, // write reply to http/ws
	handlerFunc func(
		*MessagingClient,
		IMessagesRegistry,
		json.RawMessage,
		string,
		string,
		bool,
		*mongo.Database,
		interface{},
	),
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

			handlerFunc(
				MessagingClient,
				MessagesRegistryClient,
				d.Body,
				"",
				"",
				true,
				dbRef,
				writeReplyTo,
			)
		} else {
			log.Debug("HandleRefreshEvent: correlationId not found!")
			log.Debugf("HandleRefreshEvent: Received ReplyTo: %s", d.ReplyTo)

			// new event sent by Outside App
			// handle and respond back to Outside App
			handlerFunc(
				MessagingClient,
				MessagesRegistryClient,
				d.Body,
				d.ReplyTo,
				correlationId,
				false,
				dbRef,
				writeReplyTo,
			)
		}
	}
}
