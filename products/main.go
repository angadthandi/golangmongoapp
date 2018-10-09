package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/angadthandi/golangmongoapp/products/config"
	"github.com/angadthandi/golangmongoapp/products/messages"
)

// send message on rabbitmq
func sendMQ(
	w http.ResponseWriter,
	r *http.Request,
) {
	messages.Send()
	fmt.Fprintf(w, "ProductApp Send Page! %s", "send")
}

// receive message on rabbitmq
func receiveMQ(
	w http.ResponseWriter,
	r *http.Request,
) {
	messages.Receive()
	fmt.Fprintf(w, "ProductApp Receive Page! %s", "receive")
}

// handler for home
func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Products Home Page! %s", r.URL.Path[1:])
}

func main() {
	// start receiver
	go messages.Receive()

	fmt.Printf("Listening on Port: %v", config.ServerPort)

	http.HandleFunc("/", home)
	http.HandleFunc("/send", sendMQ)
	http.HandleFunc("/receive", receiveMQ)

	// start http web server
	log.Fatal(http.ListenAndServe(config.ServerPort, nil))
}
