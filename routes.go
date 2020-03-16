package main

import "github.com/NOVAPokemon/utils"

const GeneratePokemonName = "GENERATE_POKEMON"
const GenerateItemName = "GENERATE_ITEM"

const GET = "GET"
const POST = "POST"

var routes = utils.Routes{
	utils.Route{
		Name:        GeneratePokemonName,
		Method:      GET,
		Pattern:     "/generate_pokemon",
		HandlerFunc: GeneratePokemon,
	},
	utils.Route{
		Name:        GenerateItemName,
		Method:      GET,
		Pattern:     "/generate_item",
		HandlerFunc: GenerateItem,
	},
}
