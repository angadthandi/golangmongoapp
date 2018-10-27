package main

import (
	"fmt"
	"net/http"

	"github.com/angadthandi/golangmongoapp/products/config"
	"github.com/angadthandi/golangmongoapp/products/messages"
)

var MessagingClient messages.IMessagingClient

// send message on rabbitmq
func sendMQ(
	w http.ResponseWriter,
	r *http.Request,
) {
	msg := "Product Info!"

	MessagingClient.Send(
		config.ExchangeName,
		config.ExchangeType,
		config.ProductsPublishRoutingKey,
		[]byte(msg),
	)
	fmt.Fprintf(w, "ProductApp Send Page! %s", "send")
}

// receive message on rabbitmq
func receiveMQ(
	w http.ResponseWriter,
	r *http.Request,
) {
	var sbindRoutingKeys []string
	sbindRoutingKeys = append(sbindRoutingKeys, config.GoappRoutingKey)

	MessagingClient.Receive(
		config.ExchangeName,
		config.ExchangeType,
		sbindRoutingKeys,
	)
	fmt.Fprintf(w, "ProductApp Receive Page! %s", "receive")
}

// handler for home
func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Products Home Page! %s", r.URL.Path[1:])
}

func main() {
	// initialize logging
	initLogger()
	log.Info("Starting products main...")

	// connect to RabbitMQ
	MessagingClient = &messages.MessagingClient{}
	MessagingClient.Connect()

	defer MessagingClient.Close()

	// start receiver
	// listen to messages from RabbitMQ
	// sent by other micro services
	go MessagingClient.Receive(
		config.ExchangeName,
		config.ExchangeType,
		[]string{config.GoappRoutingKey},
	)

	fmt.Printf("Listening on Port: %v", config.ServerPort)

	http.HandleFunc("/", home)
	http.HandleFunc("/send", sendMQ)
	http.HandleFunc("/receive", receiveMQ)

	// start http web server
	log.Fatal(http.ListenAndServe(config.ServerPort, nil))
}
