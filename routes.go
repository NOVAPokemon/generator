package main

import "github.com/NOVAPokemon/utils"

const GeneratePokemonName = "GENERATE_POKEMON"

const GET = "GET"
const POST = "POST"

var routes = utils.Routes{
	utils.Route{
		Name:        GeneratePokemonName,
		Method:      GET,
		Pattern:     "/generate_pokemon",
		HandlerFunc: GeneratePokemon,
	},
}
