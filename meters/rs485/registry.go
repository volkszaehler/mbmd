package rs485

import (
	"log"
	"strings"
)

// Producers is the registry of Producer factory functions
var Producers = make(map[string]func() Producer)

// Register registers a producer implementation
func Register(factory func() Producer) {
	p := factory()
	meterType := strings.ToUpper(p.Type())

	if _, ok := Producers[meterType]; ok {
		log.Fatalf("Cannot register duplicate meter type %s", meterType)
	}

	Producers[meterType] = factory
}
