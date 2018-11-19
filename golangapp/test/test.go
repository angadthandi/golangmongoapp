package test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/angadthandi/golangmongoapp/golangapp/config"
	"github.com/angadthandi/golangmongoapp/golangapp/jsondefinitions"
	"github.com/angadthandi/golangmongoapp/golangapp/messages"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	log "github.com/sirupsen/logrus"
)

type Person struct {
	// OID   objectid.ObjectID `bson:"_id"`
	Name  string
	Phone string
}

// Test MongoDB Insert/Read
// func Echo(db *mgo.Session) string {

// // connect to database
// dbSession, err := mgo.Dial(
// 	"mongodb://" +
// 		config.MongoDBUsername + ":" +
// 		config.MongoDBPassword + "@" +
// 		config.MongoDBServiceName +
// 		config.MongoDBPort)
// if err != nil {
// 	log.Fatalf("mongodb connection error : %v", err)
// }

// defer dbSession.Close()

// dbSession.SetMode(mgo.Monotonic, true)

// 	c := db.DB("testdb").C("people")
// 	_ = c.Insert(
// 		&Person{"Alex", "+55 89 9556 9111"},
// 		&Person{"John", "+55 98 8402 3256"})

// 	result := Person{}
// 	_ = c.Find(bson.M{"name": "Alex"}).One(&result)

// 	log.Debugf("Echo Response: %v", result.Phone)
// 	return result.Phone
// }

// https://medium.com/@wembleyleach/how-to-use-the-official-mongodb-go-driver-9f8aff716fdb
// https://godoc.org/github.com/mongodb/mongo-go-driver/mongo
// Test MongoDB Insert/Read with mongo-go-driver/mongo
func Echo(dbClient *mongo.Client) string {
	c := dbClient.Database("testdb").Collection("people")
	// p := Person{
	// 	// OID:   objectid.New(),
	// 	Name:  "Alex",
	// 	Phone: "+55 89 9556 9111",
	// }
	// _, err := c.InsertOne(
	// 	context.Background(),
	// 	p,
	// )

	_, err := c.InsertOne(
		context.Background(),
		bson.D{
			{Key: "Name", Value: "Alex"},
			{Key: "Phone", Value: "+55 89 9556 9111"},
		},
	)

	if err != nil {
		log.Errorf("Echo Collection Insert error: %v", err)
	}

	// ret := bson.NewDocument()
	// filter := bson.NewDocument(bson.EC.String("Name", "Alex"))
	// err = c.FindOne(context.Background(), filter).Decode(ret)
	ret := bson.D{}
	filter := bson.D{
		{Key: "Name", Value: "Alex"},
	}
	err = c.FindOne(context.Background(), filter).Decode(&ret)
	if err != nil {
		log.Errorf("Echo Document Find error: %v", err)
	}

	log.Debugf("Echo Document decoded Result: %v", ret)
	var person Person

	// person.Name = ret.Lookup("Name").StringValue()
	// person.Phone = ret.Lookup("Phone").StringValue()

	retMap := ret.Map()
	log.Debugf("Echo retMap: %v", retMap)

	name, ok := retMap["Name"].(string)
	if ok {
		person.Name = name
	}
	phone, ok := retMap["Phone"].(string)
	if ok {
		person.Phone = phone
	}
	log.Debugf("Echo person.Name: %v", person.Name)
	log.Debugf("Echo person.Phone: %v", person.Phone)

	return person.Phone
}

// Test Routes ----------------------------------------

// send message on rabbitmq
func SendMQ(
	w http.ResponseWriter,
	r *http.Request,
	MessagingClient messages.IMessagingClient,
	MessagesRegistryClient messages.IMessagesRegistry,
) {
	// var m struct{ Data string }
	// m.Data = "GoApp Publish Message!"
	var m jsondefinitions.GenericMessageSend
	m.Type = "GetProducts"

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
		false,
	)

	log.Debugf("GolangApp Send Page! %s", "send")
	fmt.Fprintf(w, "GolangApp Send Page! %s", "send")
}

// test handler
func TestHandler(
	w http.ResponseWriter,
	r *http.Request,
	dbClient *mongo.Client,
) {
	log.Debugf("Test Page! %s", Echo(dbClient))
	fmt.Fprintf(w, "Test Page! %s", Echo(dbClient))
}
