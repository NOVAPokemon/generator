package main

type GeneratorServerConfig struct {
	IntervalBetweenGenerations int `json:"interval_generate"` //in minutes
	NumberOfPokemonsToGenerate int `json:"pokemons_to_generate"`

	MaxLevel  float64 `json:"max_level"`
	MaxHP     float64 `json:"max_hp"`
	MaxDamage float64 `json:"max_damage"`
}
