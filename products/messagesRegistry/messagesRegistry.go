package messagesRegistry

import log "github.com/sirupsen/logrus"

type CorrelationMapData struct {
	SentToAppName  string
	SentToAppEvent string
}

type TCorrelationMap map[string]CorrelationMapData

var CorrelationMap TCorrelationMap

// Defines our interface for registring
// sent messages data with correlationId
type IMessagesRegistry interface {
	InitCorrelationMap()
	GetCorrelationData(
		correlationId string,
	) (sendingAppNameSendingAppEvent CorrelationMapData, exists bool)
	SetCorrelationMapData(
		correlationId string,
		sentToAppName string,
		sentToAppEvent string,
	)
	DeleteCorrelationMapData(correlationId string)
}

// Real implementation, encapsulates a pointer to CorrelationMap
type MessagesRegistryClient struct {
	correlationMap TCorrelationMap
}

func (m *MessagesRegistryClient) InitCorrelationMap() {
	m.correlationMap = make(TCorrelationMap)
	log.Debugf("messages_registry: init CorrelationMap: %v",
		m.correlationMap)
}

// returns the struct
// SentToAppName
// SentToAppEvent
func (m *MessagesRegistryClient) GetCorrelationData(
	correlationId string,
) (CorrelationMapData, bool) {
	log.Debug("messages_registry: GetCorrelationData")
	if _, ok := m.correlationMap[correlationId]; !ok {
		log.Debugf(`messages_registry: GetCorrelationData
		No match found for correlationId: %v`,
			correlationId)
		return CorrelationMapData{}, false
	}

	log.Debugf("messages_registry: GetCorrelationData for correlationId: %v",
		correlationId)
	log.Debugf("messages_registry: GetCorrelationData CorrelationMap: %v",
		m.correlationMap)
	return m.correlationMap[correlationId], true
}

func (m *MessagesRegistryClient) SetCorrelationMapData(
	correlationId string,
	sentToAppName string,
	sentToAppEvent string,
) {
	log.Debugf(`messages_registry: SetCorrelationData:
	correlationId: %v, SentToAppName: %s, SentToAppEvent: %s`,
		correlationId, sentToAppName, sentToAppEvent)

	m.correlationMap[correlationId] = CorrelationMapData{
		SentToAppName:  sentToAppName,
		SentToAppEvent: sentToAppEvent,
	}
	log.Debugf("messages_registry: SetCorrelationMapData CorrelationMap: %v",
		m.correlationMap)
}

func (m *MessagesRegistryClient) DeleteCorrelationMapData(correlationId string) {
	log.Debugf(`messages_registry: DeleteCorrelationMapData:
	correlationId: %v`, correlationId)

	delete(m.correlationMap, correlationId)
	log.Debugf("messages_registry: DeleteCorrelationMapData CorrelationMap: %v",
		m.correlationMap)
}
