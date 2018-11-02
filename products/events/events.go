package events

import (
	"encoding/json"

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

func HandleRefreshEvent(d amqp.Delivery) {
	body := d.Body
	consumerTag := d.ConsumerTag
	correlationId := d.CorrelationId
	updateToken := &UpdateToken{}
	err := json.Unmarshal(body, updateToken)
	if err != nil {
		log.Printf("Problem parsing UpdateToken: %v", err.Error())
	} else {
		log.Debugf("HandleRefreshEvent: Received a CorrelationId: %s", correlationId)
		log.Debugf("HandleRefreshEvent: Received a ConsumerTag: %s", consumerTag)
		log.Debugf("HandleRefreshEvent: Received a message: %s", body)

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
