package main

import (
	"github.com/NOVAPokemon/utils"
	"github.com/NOVAPokemon/utils/api"
)

const GET = "GET"

const CatchWildPokemonName = "CATCH_WILD_POKEMON"
const GenerateRaidBossName = "GENERATE_RAID_BOSS"

var routes = utils.Routes{
	utils.Route{
		Name:        CatchWildPokemonName,
		Method:      GET,
		Pattern:     api.CatchWildPokemonRoute,
		HandlerFunc: HandleCatchWildPokemon,
	},
	utils.Route{
		Name:        GenerateRaidBossName,
		Method:      GET,
		Pattern:     api.GenerateRaidBossRoute,
		HandlerFunc: HandleGenerateRaidBoss,
	},
}
