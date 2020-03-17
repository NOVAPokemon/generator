package main

import (
	"encoding/json"
	"github.com/NOVAPokemon/utils"
	"github.com/NOVAPokemon/utils/database/pokemon"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"math/rand"
	"time"
)

const PokemonsFile = "pokemons.json"
const MaxLevel = 100
const MaxHP = 500
const MaxDamage = 250

const numberOfPokemonsToGenerate = 10
const intervalBetweenGenerations = 2 * time.Minute

var pokemonSpecies = loadPokemons()

func main() {
	rand.Seed(time.Now().Unix())

	log.Infof("Starting GENERATOR server...\n")

	for {
		log.Info("Refreshing wild pokemons...")
		cleanWildPokemons()
		generateWildPokemons()
		time.Sleep(intervalBetweenGenerations)
	}
}

func cleanWildPokemons() {
	err := pokemon.DeleteWildPokemons()

	if err != nil {
		return
	}
}

func generateWildPokemons() {
	for i := 0; i < numberOfPokemonsToGenerate; i++ {
		toAdd := utils.Pokemon{
			Id:      primitive.NewObjectID(),
			Owner:   primitive.NewObjectID(),
			Species: pokemonSpecies[rand.Intn(len(pokemonSpecies))],
			Level:   rand.Intn(MaxLevel),
			HP:      rand.Intn(MaxHP),
			Damage:  rand.Intn(MaxDamage),
		}

		err, _ := pokemon.AddWildPokemon(toAdd)

		if err != nil {
			log.Error("Error adding wild pokemon")
			log.Error(err)
		}
	}
}

func getNumberOfItems() int {
	// TODO change this for mongoDB call to size of collection
	return 20
}

func loadPokemons() []string {
	data, err := ioutil.ReadFile(PokemonsFile)
	if err != nil {
		log.Fatal("Error loading pokemons file")
		return nil
	}

	var pokemonNames []string
	err = json.Unmarshal(data, &pokemonNames)

	if err != nil {
		log.Errorf("Error unmarshalling pokemons name")
		log.Fatal(err)
	}

	return pokemonNames
}
