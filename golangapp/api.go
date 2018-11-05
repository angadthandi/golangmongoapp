package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/angadthandi/golangmongoapp/golangapp/jsondefinitions"
	"github.com/angadthandi/golangmongoapp/golangapp/messages"
	"github.com/angadthandi/golangmongoapp/golangapp/messagesRegistry"
	"github.com/mongodb/mongo-go-driver/mongo"
	log "github.com/sirupsen/logrus"
)

// handler for API
func api(
	w http.ResponseWriter,
	r *http.Request,
	dbClient *mongo.Client,
	MessagingClient messages.IMessagingClient,
	MessagesRegistryClient messagesRegistry.IMessagesRegistry,
) {
	var resp jsondefinitions.GenericAPIResponse

	resp.Api = "Test"
	msg := "Test JSON Message!"
	resp.Message = msg

	b, err := json.Marshal(resp)
	if err != nil {
		log.Errorf("API JSON Marshal error: %v", err)
		return
	}

	log.Debugf("API JSON Response: %v", string(b))
	fmt.Fprintf(w, "API JSON Response: %v", string(b))
}

// handler for home
func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home Page! %s", r.URL.Path[1:])
}
