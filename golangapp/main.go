package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/angadthandi/golangmongoapp/golangapp/test"
)

var ServerPort = ":9002"

func testHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Test Page! %s", test.Echo())
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home Page! %s", r.URL.Path[1:])
}

func main() {

	fmt.Printf("Listening on Port: %v", ServerPort)

	// start http web server
	http.HandleFunc("/test", testHandler)
	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(ServerPort, nil))
}
