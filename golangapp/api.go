package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/angadthandi/golangmongoapp/golangapp/jsondefinitions"
	"github.com/angadthandi/golangmongoapp/golangapp/messages"
	log "github.com/sirupsen/logrus"
	mgo "gopkg.in/mgo.v2"
)

// handler for API
func api(
	w http.ResponseWriter,
	r *http.Request,
	db *mgo.Session,
	MessagingClient messages.IMessagingClient,
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
