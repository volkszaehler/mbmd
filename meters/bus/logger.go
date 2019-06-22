package bus

import "log"

// Logger is an injectable logger for grid-x modbus implementation
type Logger interface {
	Printf(format string, v ...interface{})
}

// LogLogger implements the modbus logger for standard log output
type LogLogger struct{}

func (l *LogLogger) Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}
