package rs485

import . "github.com/volkszaehler/mbmd/meters"

func init() {
	Register("WS100", NewWS100Producer)
}

type WS100Producer struct {
	Opcodes
}

func NewWS100Producer() Producer {
	/**
	 * Opcodes as defined by B+G e-tech WS100.
	 * See https://data.stromzähler.eu/manuals/bg_ws100serie_de.pdf
	 */
	ops := Opcodes{
		Voltage:       0x0100,
		Current:       0x0102,
		Power:         0x0104,
		ApparentPower: 0x0106,
		ReactivePower: 0x0108,
		Import:        0x010e,
		Export:        0x0118,
		Sum:           0x0122,
		Frequency:     0x010a,
		Cosphi:        0x010b,
	}
	return &WS100Producer{Opcodes: ops}
}

func (p *WS100Producer) Description() string {
	return "B+G e-tech WS100"
}

func (p *WS100Producer) snip(iec Measurement, readlen uint16, sign signedness, transform RTUTransform, scaler ...float64) Operation {
        snip := Operation{
                FuncCode:  ReadHoldingReg,
                OpCode:    p.Opcodes[iec],
                ReadLen:   readlen,
                Transform: transform,
                IEC61850:  iec,
        }

        if len(scaler) > 0 {
                snip.Transform = MakeScaledTransform(snip.Transform, scaler[0])
        }

        return snip
}

// snip16u creates modbus operation for single register
func (p *WS100Producer) snip16u(iec Measurement, scaler ...float64) Operation {
        return p.snip(iec, 1, unsigned, RTUUint16ToFloat64, scaler...)
}

// snip32u creates modbus operation for double register
func (p *WS100Producer) snip32u(iec Measurement, scaler ...float64) Operation {
        return p.snip(iec, 2, unsigned, RTUUint32ToFloat64, scaler...)
}

func (p *WS100Producer) Probe() Operation {
	return p.snip32u(Voltage, 1000)
}

func (p *WS100Producer) Produce() (res []Operation) {
        for _, op := range []Measurement{
                Voltage, Current,
        } {
                res = append(res, p.snip32u(op, 1000))
        }

        for _, op := range []Measurement{
                Power, ApparentPower, ReactivePower,
        } {
                res = append(res, p.snip32u(op, 1))
        }

        for _, op := range []Measurement{
                Import, Export, Sum,
        } {
                res = append(res, p.snip32u(op, 100))
        }

        for _, op := range []Measurement{
                Frequency,
        } {
                res = append(res, p.snip16u(op, 10))
        }

        for _, op := range []Measurement{
                Cosphi,
        } {
                res = append(res, p.snip16u(op, 1000))
        }

	return res
}
