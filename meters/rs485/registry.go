package rs485

import (
	"log"
	"strings"
)

var producers = make(map[string]func() Producer)

// Register registers a producer implementation
func Register(factory func() Producer) {
	p := factory()
	meterType := strings.ToUpper(p.Type())

	if _, ok := producers[meterType]; ok {
		log.Fatalf("Cannot register duplicate meter type %s", meterType)
	}

	producers[meterType] = factory
}
