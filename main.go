package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const host = "localhost"
const Port = 8002

func main() {
	router := NewRouter()
	addr := fmt.Sprintf("%s:%d", host, Port)

	log.Info("Starting AUTHENTICATION server...")
	log.Fatal(http.ListenAndServe(addr, router))
}
