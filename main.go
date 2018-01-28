package main

import (
	"log"
	"net/http"
	"fmt"
	"io/ioutil"
)

var webhookPath = fmt.Sprintf("/%s", telegram.Token)

func main() {
	// TODO Create a config file
	dbConfig := DbConfig{
		"0.0.0.0",
		5432,
		"dennis",
		"dennis_test",
		"dennis",
		"disable",
	}

	// Set up DB
	dbConfig.Open()

	http.HandleFunc("/healthcheck", healthcheck)
	http.HandleFunc(webhookPath, webhook)

	log.Printf("main: starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func healthcheck(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("ok"))
}

func webhook(w http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	log.Printf("main: incoming message - %s", body)
	converse(body)
	w.Write([]byte("received"))
}
