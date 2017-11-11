package main

import (
	"log"
	"net/http"
	"fmt"
)

var webhookPath = fmt.Sprintf("/%s", telegram.Token)

func main() {
	http.HandleFunc("/healthcheck", healthcheck)
	http.HandleFunc(webhookPath, telegram.Webhook)

	log.Printf("bot: starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func healthcheck(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("ok"))
}
