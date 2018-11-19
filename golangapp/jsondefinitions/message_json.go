package jsondefinitions

import (
	"encoding/json"
)

type GenericMessageRecieve struct {
	Type    string          `json:"type"`
	Message json.RawMessage `json:"message"`
}

type GenericMessageSend struct {
	Type    string      `json:"type"`
	Message interface{} `json:"message"`
}

type GenericErrorMessageSend struct {
	Errormessage interface{} `json:"errormessage"`
}

// MicroServiceResponseMsgForHub
// received msg with correlationId
// from other micro service
// to be sent on hub channel
type MicroServiceResponseMsgForHub struct {
	CorrelationId   string
	ReceivedJsonMsg json.RawMessage
}
