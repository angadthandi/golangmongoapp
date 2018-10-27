package main

import (
	"net/http"

	"github.com/angadthandi/golangmongoapp/golangapp/test"
	mgo "gopkg.in/mgo.v2"
)

// configure API Routes
func configureRoutes(dbSession *mgo.Session) {
	http.HandleFunc("/", home)
	http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		api(w, r, dbSession, MessagingClient)
	})

	// Test Routes --------------------------------

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		test.TestHandler(w, r, dbSession)
	})
	http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		test.SendMQ(w, r, MessagingClient)
	})
}
