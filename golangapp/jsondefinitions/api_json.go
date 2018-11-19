package jsondefinitions

import (
	"encoding/json"
)

type GenericAPIRecieve struct {
	Api     string          `json:"api"`
	Message json.RawMessage `json:"message"`
}

type GenericAPIResponse struct {
	Api     string      `json:"api"`
	Message interface{} `json:"message"`
}

// ClientUUIDCorrelationID
// used to set correlationId
// for client uuid
// to setClientCorrelationId
type ClientUUIDCorrelationID struct {
	ClientUUID          string
	ClientCorrelationId string
}
