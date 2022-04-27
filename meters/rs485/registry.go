package rs485

import (
	"log"
)

// Producers is the registry of Producer factory functions
var Producers = make(map[string]func() Producer)

// Register registers a producer implementation
func Register(typ string, factory func() Producer) {
	if _, ok := Producers[typ]; ok {
		log.Fatalf("cannot register duplicate meter type: %s", typ)
	}

	Producers[typ] = factory
}
