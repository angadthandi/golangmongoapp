package main

import (
	"net/http"

	"github.com/angadthandi/golangmongoapp/golangapp/test"
	"github.com/gorilla/mux"
	mgo "gopkg.in/mgo.v2"
)

// configure API Routes
func configureRoutes(dbSession *mgo.Session) {
	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/", home)
	r.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		api(w, r, dbSession, MessagingClient, MessagesRegistryClient)
	})

	// Test Routes --------------------------------

	r.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		test.TestHandler(w, r, dbSession)
	})
	r.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		test.SendMQ(w, r, MessagingClient, MessagesRegistryClient)
	})

	http.Handle("/", r)
}
