package bus

import "time"

type Bus interface {
	String() string
	Logger(l Logger)
	Timeout(timeout time.Duration) time.Duration
	Slave(deviceID uint8)
	Reconnect()
}
