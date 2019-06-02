package impl

import . "github.com/volkszaehler/mbmd/meters"

type RS485Core struct {
	Opcodes
}

func (p *RS485Core) ConnectionType() ConnectionType {
	return RS485
}
