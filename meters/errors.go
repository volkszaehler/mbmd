package meters

import "errors"

var (
	// ErrNaN indicates a NaN reading result
	ErrNaN = errors.New("NaN value")

	// ErrPartiallyOpened indicates a partially opened device
	ErrPartiallyOpened = errors.New("Device partially opened")
)
