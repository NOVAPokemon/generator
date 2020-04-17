package main

import (
	"encoding/json"
	"fmt"
	"github.com/NOVAPokemon/utils"
	generatordb "github.com/NOVAPokemon/utils/database/generator"
	"github.com/NOVAPokemon/utils/items"
	"github.com/NOVAPokemon/utils/pokemons"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

// Pokemons taken from https://raw.githubusercontent.com/sindresorhus/pokemon/master/data/en.json
const PokemonsFile = "pokemons.json"
const numberOfPokemonsToGenerate = 10

const MaxLevel = 100
const MaxHP = 500
const MaxDamage = 500

const stdHPDeviation = float64(MaxHP) / 20
const stdDamageDeviation = float64(MaxDamage) / 20

const ItemsFile = "items.json"
const numberOfItemsToGenerate = 20

const intervalBetweenGenerations = 2 * time.Minute

var pokemonSpecies = loadPokemons()
var itemNames = loadItems()

const host = utils.Host
const port = utils.GeneratorPort

var addr = fmt.Sprintf("%s:%d", host, port)

func main() {
	rand.Seed(time.Now().Unix())

	log.Infof("Starting GENERATOR server...\n")

	go generate()

	r := utils.NewRouter(routes)
	log.Infof("Starting STORE server in port %d...\n", port)
	log.Fatal(http.ListenAndServe(addr, r))
}

func generate() {
	for {
		log.Info("Refreshing wild pokemons...")
		cleanWildPokemons()
		generateWildPokemons(numberOfPokemonsToGenerate)
		log.Info("Refreshing catchable items...")
		cleanItems()
		generateItems()
		time.Sleep(intervalBetweenGenerations)
	}
}

func cleanWildPokemons() {
	err := generatordb.DeleteWildPokemons()

	if err != nil {
		return
	}
}

func getOneWildPokemon() *pokemons.Pokemon {
	var level, hp, damage int
	level = rand.Intn(MaxLevel-1) + 1
	log.Println("Level: ", level)
	randNormal := rand.NormFloat64()*
		(stdHPDeviation*(float64(level)/MaxLevel)) +
		(MaxHP * (float64(level) / MaxLevel))
	hp = int(randNormal)
	log.Println("HP: ", hp)

	//safeguards
	if hp < 1 {
		hp = 1
	}

	randNormal = rand.NormFloat64()*
		(stdDamageDeviation*(float64(level)/MaxLevel)) +
		(MaxDamage * (float64(level) / MaxLevel))

	damage = int(randNormal)
	log.Println("Damage: ", damage)

	//safeguards
	if damage < 1 {
		damage = 1
	}

	wildPokemon := &pokemons.Pokemon{
		Id:      primitive.NewObjectID(),
		Species: pokemonSpecies[rand.Intn(len(pokemonSpecies))],
		Level:   level,
		HP:      hp,
		MaxHP:   hp,
		Damage:  damage,
	}
	return wildPokemon
}

func generateWildPokemons(numberOfPokemonsToGenerate int) {

	for i := 0; i < numberOfPokemonsToGenerate; i++ {

		err, _ := generatordb.AddWildPokemon(*getOneWildPokemon())

		if err != nil {
			log.Error("Error adding wild pokemon")
			log.Error(err)
		}
	}
}

func generateRaidBoss() *pokemons.Pokemon { // TODO look at this
	generated := getOneWildPokemon()
	generated.Level *= 5
	generated.HP *= 5
	generated.MaxHP *= 5
	return generated
}

func cleanItems() {
	err := generatordb.DeleteCatchableItems()

	if err != nil {
		return
	}
}

func generateItems() {
	var toAdd items.Item
	for i := 0; i < numberOfItemsToGenerate; i++ {
		toAdd = items.Item{
			Id:   primitive.NewObjectID(),
			Name: itemNames[rand.Intn(len(itemNames))],
		}

		err, _ := generatordb.AddCatchableItem(toAdd)

		if err != nil {
			log.Error("Error adding item")
			log.Error(err)
		}
	}
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

	log.Infof("Loaded %d pokemon species.", len(pokemonNames))

	return pokemonNames
}

func loadItems() []string {
	data, err := ioutil.ReadFile(ItemsFile)
	if err != nil {
		log.Fatal("Error loading items file")
		return nil
	}

	var itemNames []string
	err = json.Unmarshal(data, &itemNames)

	if err != nil {
		log.Errorf("Error unmarshalling item names")
		log.Fatal(err)
	}

	log.Infof("Loaded %d items.", len(itemNames))

	return itemNames
}
