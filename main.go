package main

import (
	"encoding/json"
	"fmt"
	"github.com/NOVAPokemon/utils"
	generatordb "github.com/NOVAPokemon/utils/database/generator"
	"github.com/NOVAPokemon/utils/pokemons"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

const configFilename = "configs.json"

// Pokemons taken from https://raw.githubusercontent.com/sindresorhus/pokemon/master/data/en.json
const PokemonsFile = "pokemons.json"
const numberOfPokemonsToGenerate = 10

var (
	config = loadConfig()

	pokemonSpecies = loadPokemons()

	stdHPDeviation     = config.MaxHP / 20
	stdDamageDeviation = config.MaxDamage / 20
)

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
		time.Sleep(time.Duration(config.IntervalBetweenGenerations) * time.Minute)
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
	level = rand.Intn(int(config.MaxLevel-1)) + 1
	log.Println("Level: ", level)
	randNormal := rand.NormFloat64()*
		(stdHPDeviation*(float64(level)/config.MaxLevel)) +
		(config.MaxHP * (float64(level) / config.MaxLevel))
	hp = int(randNormal)
	log.Println("HP: ", hp)

	//safeguards
	if hp < 1 {
		hp = 1
	}

	randNormal = rand.NormFloat64()*
		(stdDamageDeviation*(float64(level)/config.MaxLevel)) +
		(config.MaxDamage * (float64(level) / config.MaxLevel))

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

func loadConfig() *GeneratorServerConfig {
	data, err := ioutil.ReadFile(configFilename)
	if err != nil {
		log.Fatal("Error loading configs")
		return nil
	}

	var config GeneratorServerConfig
	err = json.Unmarshal(data, &config)

	if err != nil {
		log.Errorf("Error unmarshalling configs")
		log.Fatal(err)
	}

	return &config
}
