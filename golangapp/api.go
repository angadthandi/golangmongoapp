package main

import (
	"fmt"
	"net/http"

	"github.com/angadthandi/golangmongoapp/golangapp/messages"
	mgo "gopkg.in/mgo.v2"
)

// handler for API
func api(
	w http.ResponseWriter,
	r *http.Request,
	db *mgo.Session,
	MessagingClient messages.IMessagingClient,
) {
	log.Printf("API Handler Page! %s", r.URL.Path[1:])
	fmt.Fprintf(w, "API Handler Page! %s", r.URL.Path[1:])
}

// handler for home
func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home Page! %s", r.URL.Path[1:])
}
