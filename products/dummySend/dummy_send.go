package dummySend

import (
	"encoding/json"

	"github.com/angadthandi/golangmongoapp/products/config"
	"github.com/angadthandi/golangmongoapp/products/messages"
	"github.com/angadthandi/golangmongoapp/products/messagesRegistry"
	log "github.com/sirupsen/logrus"
)

func DummySendToGoapp(
	MessagingClient *messages.MessagingClient,
	MessagesRegistryClient messagesRegistry.IMessagesRegistry,
	sendToRoutingKey string,
	receivedCorrelationId string,
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
	)

	log.Debugf("Products Send Page! %s", "send")
}
