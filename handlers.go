package main

import (
	"encoding/json"
	"github.com/NOVAPokemon/utils"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
)

func GeneratePokemon(w http.ResponseWriter, r *http.Request) {
	randId := rand.Intn(NumberOfPokemons)

	// TODO change for mongoDB call for a pokemon with number /randid
	pokemon := randId
	pokemonJSON, err := json.Marshal(pokemon)

	if err != nil {
		utils.HandleJSONEncodeError(&w, GeneratePokemonName, err)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(pokemonJSON)
	log.Infof("%s: %d\n", GeneratePokemonName, pokemon)
}

func GenerateItem(w http.ResponseWriter, r *http.Request) {
	randId := rand.Intn(NumberOfItems)

	// TODO change for mongoDB call for a pokemon with number /randid
	item := randId
	itemJSON, err := json.Marshal(item)

	if err != nil {
		utils.HandleJSONEncodeError(&w, GeneratePokemonName, err)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(itemJSON)
	log.Infof("%s: %d\n", GenerateItemName, item)
}
