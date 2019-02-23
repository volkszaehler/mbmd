package meters

type SunSpecDeviceDescriptor struct {
	Manufacturer string
	Model        string
	Options      string
	Version      string
	Serial       string
}

type SunSpecProducer interface {
	GetSunSpecCommonBlock() Operation
	DecodeSunSpecCommonBlock(b []byte) (SunSpecDeviceDescriptor, error)
}
