package rs485

import (
	"log"
	"strings"
)

type producerFactory func() Producer

var producers = make(map[string]producerFactory)

func Register(factory producerFactory) {
	p := factory()
	meterType := strings.ToUpper(p.Type())

	if _, ok := producers[meterType]; ok {
		log.Fatalf("Cannot register duplicate meter type %s", meterType)
	}

	producers[meterType] = factory
}
