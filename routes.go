package main

import (
	"github.com/NOVAPokemon/utils"
	"github.com/NOVAPokemon/utils/api"
)

const GET = "GET"

const CatchWildPokemonName = "CATCH_WILD_POKEMON"

var routes = utils.Routes{
	utils.Route{
		Name:        CatchWildPokemonName,
		Method:      GET,
		Pattern:     api.CatchWildPokemonRoute,
		HandlerFunc: HandleCatchWildPokemon,
	},
}
