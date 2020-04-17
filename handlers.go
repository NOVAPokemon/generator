package main

import (
	"encoding/json"
	"fmt"
	"github.com/NOVAPokemon/utils"
	"github.com/NOVAPokemon/utils/clients"
	generatordb "github.com/NOVAPokemon/utils/database/generator"
	"github.com/NOVAPokemon/utils/items"
	"github.com/NOVAPokemon/utils/tokens"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
)

const (
	maxCatchingProbability = 100
)

var httpClient = &http.Client{}

func HandleCatchWildPokemon(w http.ResponseWriter, r *http.Request) {
	authToken, err := tokens.ExtractAndVerifyAuthToken(r.Header)
	if err != nil {
		log.Error("no auth token")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var pokeball items.Item
	err = json.NewDecoder(r.Body).Decode(&pokeball)
	if err != nil {
		log.Error(err)
		return
	}

	if !pokeball.IsPokeBall() {
		log.Error("invalid item to catch")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	wildPokemons := generatordb.GetWildPokemons()
	selectedPokemon := &wildPokemons[rand.Intn(len(wildPokemons))]

	var catchingProbability float64
	if pokeball.Effect.Value == maxCatchingProbability {
		catchingProbability = 1
	} else {
		catchingProbability = 1 - ((float64(selectedPokemon.Level) / MaxLevel) * (float64(pokeball.Effect.Value) / maxCatchingProbability))
	}

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
		}
		return
	}

	log.Info(authToken.Username, " caught: ", caught)
	var trainersClient = clients.NewTrainersClient(fmt.Sprintf("%s:%d", utils.Host, utils.TrainersPort), httpClient)
	_, err = trainersClient.AddPokemonToTrainer(authToken.Username, *selectedPokemon)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	pokemonTokens := make([]string, 0, len(trainersClient.PokemonTokens))
	for _, tokenString := range trainersClient.PokemonTokens {
		pokemonTokens = append(pokemonTokens, tokenString)
	}

	w.Header()[tokens.PokemonsTokenHeaderName] = pokemonTokens
	_, err = w.Write(jsonBytes)
	if err != nil {
		log.Error(err)
	}
}

func HandleGenerateRaidBoss(w http.ResponseWriter, _ *http.Request) {
	raidBossMarshaled, _ := json.Marshal(generateRaidBoss())
	_, _ = w.Write(raidBossMarshaled)
}
