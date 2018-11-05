package main

import (
	"net/http"

	"github.com/mongodb/mongo-go-driver/mongo"

	"github.com/angadthandi/golangmongoapp/golangapp/test"
	"github.com/gorilla/mux"
)

// configure API Routes
func configureRoutes(dbClient *mongo.Client) {
	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/", home)
	r.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		api(w, r, dbClient, MessagingClient, MessagesRegistryClient)
	})

	// Test Routes --------------------------------

	r.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		test.TestHandler(w, r, dbClient)
	})
	r.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		test.SendMQ(w, r, MessagingClient, MessagesRegistryClient)
	})

	http.Handle("/", r)
}
