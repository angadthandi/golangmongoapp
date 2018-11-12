package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/angadthandi/golangmongoapp/golangapp/jsondefinitions"
	"github.com/angadthandi/golangmongoapp/golangapp/messages"
	"github.com/mongodb/mongo-go-driver/mongo"
	log "github.com/sirupsen/logrus"
)

// handler for rest/API
func API(
	w http.ResponseWriter,
	r *http.Request,
	dbClient *mongo.Client,
	MessagingClient messages.IMessagingClient,
	MessagesRegistryClient messages.IMessagesRegistry,
) {
	var resp jsondefinitions.GenericAPIResponse

	resp.Api = "Test"
	msg := "Test JSON Message!"
	resp.Message = msg

	b, err := json.Marshal(resp)
	if err != nil {
		log.Errorf("rest/API JSON Marshal error: %v", err)
		return
	}

	log.Debugf("rest/API JSON Response: %v", string(b))
	fmt.Fprintf(w, "rest/API JSON Response: %v", string(b))
}

// handler for home
func Home(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "Home Page! %s", r.URL.Path[1:])

	log.Debug("Home Page!")
	http.FileServer(http.Dir("./public/home.html"))
	// http.FileServer(http.Dir("./vendor/home.html"))
}
