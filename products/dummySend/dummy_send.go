package dummySend

import (
	"encoding/json"

	"github.com/angadthandi/golangmongoapp/products/config"
	"github.com/angadthandi/golangmongoapp/products/messages"
	log "github.com/sirupsen/logrus"
)

func DummySendToGoapp(
	MessagingClient *messages.MessagingClient,
	MessagesRegistryClient messages.IMessagesRegistry,
	sendToRoutingKey string,
	receivedCorrelationId string,
	isReplyMessage bool,
) {
	var m struct{ Data string }
	m.Data = "Products Publish Message!"

	b, err := json.Marshal(m)
	if err != nil {
		log.Errorf("Products: send: unable to marshal: %v", err)
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

	log.Debugf("Products Send Page! %s", "send")
}
