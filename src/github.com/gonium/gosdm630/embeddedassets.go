package sdm630

import (
	"sync"
	"time"
	"github.com/aprice/embed/loader"
)

var _embeddedContentLoader loader.Loader
var _initOnce sync.Once

// GetEmbeddedContent returns the Loader for embedded content files.
func GetEmbeddedContent() loader.Loader {
	_initOnce.Do(_initEmbeddedContent)
	return _embeddedContentLoader
}

func _initEmbeddedContent() {
	l := loader.New()

	l.Add(&loader.Content{
		Path:   "/index.tmpl",
		Hash:    "PsrhmF_YQqrRi-sw__jhCA",
		Modified: time.Unix(1497428649, 0),
		Raw: `
PCFET0NUWVBFIGh0bWw+CjxodG1sIGxhbmc9ImVuIj4KICA8aGVhZD4KICAgIDxtZXRhIGNoYXJzZXQ9
InV0Zi04Ij4KCQkgPG1ldGEgaHR0cC1lcXVpdj0icmVmcmVzaCIgY29udGVudD0ie3suUmVsb2FkSW50
ZXJ2YWx9fSIgLz4KICAgIDx0aXRsZT5Hb1NETTYzMCBvdmVydmlldyBwYWdlPC90aXRsZT4KICA8L2hl
YWQ+CiAgPGJvZHk+CgkJPHByZT4KU0RNNjMwIEhUVFAgc2VydmVyLCB2ZXJzaW9uIHt7LlNvZnR3YXJl
VmVyc2lvbn19LiBSZWxvYWRpbmcgZXZlcnkge3suUmVsb2FkSW50ZXJ2YWx9fSBzZWNvbmRzLgp7ey5D
b250ZW50fX0KCQk8L3ByZT4KICA8L2JvZHk+CjwvaHRtbD4KCg
`,
	})

	_embeddedContentLoader = l
}
