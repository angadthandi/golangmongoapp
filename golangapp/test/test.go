package test

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/angadthandi/golangmongoapp/golangapp/config"
	"github.com/angadthandi/golangmongoapp/golangapp/messages"
	"github.com/angadthandi/golangmongoapp/golangapp/messagesRegistry"
	log "github.com/sirupsen/logrus"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Person struct {
	Name  string
	Phone string
}

// Test MongoDB Insert/Read
func Echo(db *mgo.Session) string {
	c := db.DB("testdb").C("people")
	_ = c.Insert(
		&Person{"Alex", "+55 89 9556 9111"},
		&Person{"John", "+55 98 8402 3256"})

	result := Person{}
	_ = c.Find(bson.M{"name": "Alex"}).One(&result)

	log.Debugf("Echo Response: %v", result.Phone)
	return result.Phone
}

// Test Routes ----------------------------------------

// send message on rabbitmq
func SendMQ(
	w http.ResponseWriter,
	r *http.Request,
	MessagingClient messages.IMessagingClient,
	MessagesRegistryClient messagesRegistry.IMessagesRegistry,
) {
	var m struct{ Data string }
	m.Data = "GoApp Publish Message!"

	b, err := json.Marshal(m)
	if err != nil {
		log.Errorf("GoApp: send: unable to marshal: %v", err)
	}

	MessagingClient.Send(
		config.ExchangeName,
		config.ExchangeType,
		config.ProductsRoutingKey,
		config.GoappRoutingKey,
		b,
		MessagesRegistryClient,
		"",
	)

	log.Debugf("GolangApp Send Page! %s", "send")
	fmt.Fprintf(w, "GolangApp Send Page! %s", "send")
}

// test handler
func TestHandler(
	w http.ResponseWriter,
	r *http.Request,
	db *mgo.Session,
) {
	log.Debugf("Test Page! %s", Echo(db))
	fmt.Fprintf(w, "Test Page! %s", Echo(db))
}
