package main

import (
	"encoding/json"
	"github.com/NOVAPokemon/utils"
	generatordb "github.com/NOVAPokemon/utils/database/generator"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"math/rand"
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
const numberOfItemsToGenerate = 10

const intervalBetweenGenerations = 2 * time.Minute

var pokemonSpecies = loadPokemons()
var itemNames = loadItems()

func main() {
	rand.Seed(time.Now().Unix())

	log.Infof("Starting GENERATOR server...\n")

	for {
		log.Info("Refreshing wild pokemons...")
		cleanWildPokemons()
		generateWildPokemons()
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

func generateWildPokemons() {
	var level, hp, damage int

	for i := 0; i < numberOfPokemonsToGenerate; i++ {
		level = rand.Intn(MaxLevel-1) + 1
		log.Println("Level: ", level)
		randNormal := rand.NormFloat64()*
			(stdHPDeviation*(float64(level)/MaxLevel)) + +
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

		toAdd := utils.Pokemon{
			Id:      primitive.NewObjectID(),
			Species: pokemonSpecies[rand.Intn(len(pokemonSpecies))],
			Level:   level,
			HP:      hp,
			Damage:  damage,
		}

		err, _ := generatordb.AddWildPokemon(toAdd)

		if err != nil {
			log.Error("Error adding wild pokemon")
			log.Error(err)
		}
	}
}

func cleanItems() {
	err := generatordb.DeleteCatchableItems()

	if err != nil {
		return
	}
}

func generateItems() {
	for i := 0; i < numberOfItemsToGenerate; i++ {
		toAdd := utils.Item{
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
