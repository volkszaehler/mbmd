package meters

// Logger is an injectable logger for grid-x modbus implementation
type Logger interface {
	Printf(format string, v ...interface{})
}
