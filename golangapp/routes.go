package main

import (
	"net/http"

	"github.com/mongodb/mongo-go-driver/mongo"

	"github.com/angadthandi/golangmongoapp/golangapp/api/rest"
	"github.com/angadthandi/golangmongoapp/golangapp/goappsocket"
	"github.com/angadthandi/golangmongoapp/golangapp/messages"
	"github.com/angadthandi/golangmongoapp/golangapp/test"
	"github.com/gorilla/mux"
)

// configure API Routes
func configureRoutes(
	dbClient *mongo.Client,
	hub *goappsocket.Hub,
	MessagingClient messages.IMessagingClient,
	MessagesRegistryClient messages.IMessagesRegistry,
) {
	r := mux.NewRouter().StrictSlash(true)

	// r.HandleFunc("/", rest.Home)
	r.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		rest.API(w, r, dbClient, MessagingClient, MessagesRegistryClient)
	})

	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		goappsocket.ServeWs(
			hub,
			w,
			r,
			dbClient,
			MessagingClient,
			MessagesRegistryClient,
		)
	})

	// Test Routes --------------------------------
	r.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		test.TestHandler(w, r, dbClient)
	})
	r.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		test.SendMQ(w, r, MessagingClient, MessagesRegistryClient)
	})

	// // static files
	// r.HandleFunc("/vendor/", func(w http.ResponseWriter, r *http.Request) {
	// 	http.StripPrefix("/vendor/",
	// 		http.FileServer(http.Dir("./public")))
	// })
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/")))

	http.Handle("/", r)
}
