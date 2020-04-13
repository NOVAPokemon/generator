package main

import (
	"encoding/json"
	"fmt"
	"github.com/NOVAPokemon/utils"
	"github.com/NOVAPokemon/utils/clients"
	generatordb "github.com/NOVAPokemon/utils/database/generator"
	"github.com/NOVAPokemon/utils/tokens"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"math/rand"
	"net/http"
)

var httpClient = &http.Client{}

func HandleCatchWildPokemon(w http.ResponseWriter, r *http.Request) {
	authToken, err := tokens.ExtractAndVerifyAuthToken(r.Header)
	if err != nil {
		log.Error("no auth token")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	authTokenString := r.Header.Get(tokens.AuthTokenHeaderName)

	wildPokemons := generatordb.GetWildPokemons()
	selectedPokemon := &wildPokemons[rand.Intn(len(wildPokemons))]

	catchingProbability := 1 - (float64(selectedPokemon.Level) / MaxLevel)

	log.Info("catching probability: ", catchingProbability)

	caught := rand.Float64() <= catchingProbability
	caughtMessage := clients.CaughtPokemonMessage{
		Caught: caught,
	}

	jsonBytes, err := json.Marshal(caughtMessage)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !caught {
		_, err = w.Write(jsonBytes)
		if err != nil {
			log.Error(err)
			return
		}
	}

	log.Info(authToken.Username, " caught: ", caught)
	var trainersClient = clients.NewTrainersClient(fmt.Sprintf("%s:%d", utils.Host, utils.TrainersPort), httpClient)
	newPokemons, err := trainersClient.AddPokemonToTrainer(authToken.Username, *selectedPokemon)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	v := make([]string, 0, len(trainersClient.PokemonTokens))
	for _, value := range trainersClient.PokemonTokens {
		v = append(v, value)
	}

	w.Header()[tokens.PokemonsTokenHeaderName] = v
	toSend, _ := json.Marshal(newPokemons)
	_, err = w.Write(toSend)
	if err != nil {
		log.Error(err)
		return
	}
}

func loadShopItems() ([]utils.StoreItem, map[string]utils.StoreItem) {

	data, err := ioutil.ReadFile(ItemsFile)
	if err != nil {
		log.Errorf("Error loading items file ")
		log.Fatal(err)
		panic(err)
	}

	var itemsArr []utils.StoreItem
	err = json.Unmarshal(data, &itemsArr)

	var itemsMap = make(map[string]utils.StoreItem, len(itemsArr))
	for _, item := range itemsArr {
		itemsMap[item.Name] = item
	}

	if err != nil {
		log.Errorf("Error unmarshalling item names")
		log.Fatal(err)
		panic(err)
	}

	log.Infof("Loaded %d items.", len(itemsArr))

	return itemsArr, itemsMap
}
