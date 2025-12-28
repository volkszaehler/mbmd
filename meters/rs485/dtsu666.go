package rs485

import . "github.com/volkszaehler/mbmd/meters"
import "time"

func init() {
	Register("DTSU666", NewDTSU666Producer)
}

type DTSU666Producer struct {
	Opcodes
}

func NewDTSU666Producer() Producer {
	/**
	 *Protocol definition see: https://github.com/FinestArcadeArt/Sungrow_DTSU666_esphome/blob/master/PDFs/DTSU666%20Meter%20Communication%20Protocol_20210601.pdf
	 */
	ops := Opcodes{

		//L1
		VoltageL1:       0x0061,
		CurrentL1:       0x0064,
		PowerL1:         0x0164,

		//L2
		VoltageL2:       0x0062,
		CurrentL2:       0x0065,
		PowerL2:         0x0166,
		
		//L3
		VoltageL3:       0x0063,
		CurrentL3:       0x0066,
		PowerL3:         0x0168,
		
		//Total
		Frequency:       0x0077,
		Power:           0x016A,
		Import:			 0x000A,
		Export:			 0x0014,
	}
	return &DTSU666Producer{Opcodes: ops}
}

func (p *DTSU666Producer) Description() string {
	return "CHINT DTSU666"
}

func (p *DTSU666Producer) snip16u(iec Measurement, scaler ...float64) Operation {
	operation := Operation{
		FuncCode:  ReadInputReg,
		OpCode:    p.Opcode(iec),
		ReadLen:   1,
		IEC61850:  iec,
		Transform: RTUUint16ToFloat64,
	}
	
	operation.Transform = MakeScaledTransform(operation.Transform, scaler[0])
	
	return operation
}

func (p *DTSU666Producer) snip32s(iec Measurement, scaler ...float64) Operation {
	operation := Operation{
		FuncCode:  ReadInputReg,
		OpCode:    p.Opcode(iec),
		ReadLen:   2,
		IEC61850:  iec,
		Transform: RTUInt32ToFloat64,
	}
	
	operation.Transform = MakeScaledTransform(operation.Transform, scaler[0])
		
	return operation
}

func (p *DTSU666Producer) snip32u(iec Measurement, scaler ...float64) Operation {
	operation := Operation{
		FuncCode:  ReadInputReg,
		OpCode:    p.Opcode(iec),
		ReadLen:   2,
		IEC61850:  iec,
		Transform: RTUUint32ToFloat64,
	}

	operation.Transform = MakeScaledTransform(operation.Transform, scaler[0])

	return operation
}

func (p *DTSU666Producer) Probe() Operation {
	return p.snip16u(VoltageL1)
}

func (p *DTSU666Producer) Produce() (res []Operation) {
	for _, op := range []Measurement{
		VoltageL1, VoltageL2, VoltageL3, 
	} {
		res = append(res, p.snip16u(op,10))
	}

    for _, op := range []Measurement{
		CurrentL1, CurrentL2, CurrentL3, Frequency, 
	} {
		res = append(res, p.snip16u(op,100))
	}
	
	for _, op := range []Measurement{
		PowerL1, PowerL2, PowerL3, Power,
	} {
		res = append(res, p.snip32s(op,1))
    }

	for _, op := range []Measurement{
		Import, Export,
	} {
		res = append(res, p.snip32u(op,100))
		time.Sleep(50 * time.Millisecond)
	}


	return res
}
