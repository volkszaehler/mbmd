package impl

import . "github.com/gonium/gosdm630/meters"

type RS485Core struct {
	Opcodes
}

func (p *RS485Core) ConnectionType() ConnectionType {
	return RS485
}
