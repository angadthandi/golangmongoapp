package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/angadthandi/golangmongoapp/golangapp/test"
	mgo "gopkg.in/mgo.v2"
)

var (
	ServerPort           = ":9002"          // http server port
	mongodatabaseservice = "mongodbservice" // mongodb service name
)

// test handler
func testHandler(
	w http.ResponseWriter,
	r *http.Request,
	db *mgo.Session,
) {
	fmt.Fprintf(w, "Test Page! %s", test.Echo(db))
}

// handler for home
func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home Page! %s", r.URL.Path[1:])
}

func main() {
	// connect to database
	dbSession, err := mgo.Dial(mongodatabaseservice)
	if err != nil {
		log.Fatalf("mongodb connection error : %v", err)
	}

	defer dbSession.Close()

	dbSession.SetMode(mgo.Monotonic, true)

	fmt.Printf("Listening on Port: %v", ServerPort)

	// start http web server
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		testHandler(w, r, dbSession)
	})
	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(ServerPort, nil))
}
