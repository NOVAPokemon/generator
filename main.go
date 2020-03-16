package main

import (
	"fmt"
	"github.com/NOVAPokemon/utils"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	"time"
)

const host = "localhost"
const Port = 8002

var (
	NumberOfItems    int
	NumberOfPokemons int
)

func main() {
	rand.Seed(time.Now().Unix())

	NumberOfPokemons = getNumberOfPokemons()
	NumberOfItems = getNumberOfItems()

	router := utils.NewRouter(routes)
	addr := fmt.Sprintf("%s:%d", host, Port)

	log.Infof("Starting GENERATOR server in port %d...\n", Port)
	log.Fatal(http.ListenAndServe(addr, router))
}

func getNumberOfPokemons() int {
	// TODO change this for mongoDB call to size of collection
	return 20
}

func getNumberOfItems() int {
	// TODO change this for mongoDB call to size of collection
	return 20
}
