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
