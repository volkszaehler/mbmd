package loader

import (
	"net/http"
)

// Loader types expose content which may be embedded in a binary or loaded
// from disk at runtime.
type Loader interface {
	http.Handler

	// GetContents returns the contents of the file at path, or
	// os.ErrNoExist if no such file is found.
	GetContents(path string) ([]byte, error)
}
