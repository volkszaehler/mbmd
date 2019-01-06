package meters

import (
	"log"
	"strings"
)

type FactoryFunc func() Producer

var Producers = make(map[string]FactoryFunc)

func Register(factory FactoryFunc) {
	p := factory()
	meterType := strings.ToUpper(p.Type())

	if _, ok := Producers[meterType]; ok {
		log.Fatalf("Cannot register duplicate meter type %s", meterType)
	}

	Producers[meterType] = factory
}
