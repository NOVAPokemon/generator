module github.com/NOVAPokemon/generator

go 1.13

require (
	github.com/NOVAPokemon/utils v0.0.62
	github.com/sirupsen/logrus v1.4.2
	go.mongodb.org/mongo-driver v1.3.1
)

replace github.com/NOVAPokemon/utils v0.0.62 => ../utils
