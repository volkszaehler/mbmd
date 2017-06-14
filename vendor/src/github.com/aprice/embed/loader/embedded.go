package loader

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
)

type EmbeddedLoader struct {
	http.Handler
	content map[string]*Content
	dirs    map[string]embeddedDir
}

func New() *EmbeddedLoader {
	l := &EmbeddedLoader{
		content: make(map[string]*Content),
		dirs:    map[string]embeddedDir{"/": embeddedDir{name: "/"}},
	}
	l.Handler = http.FileServer(l)
	return l
}

func (l *EmbeddedLoader) Add(c *Content) {
	if len(c.Compressed) > 0 {
		c.CompressedBytes, _ = base64.RawStdEncoding.DecodeString(c.Compressed)
	} else {
		c.CompressedBytes = make([]byte, 0)
	}
	if strings.TrimSpace(c.Raw) != "" {
		c.RawBytes, _ = base64.RawStdEncoding.DecodeString(c.Raw)
	} else if len(c.CompressedBytes) > 0 {
		gzr, _ := gzip.NewReader(bytes.NewReader(c.CompressedBytes))
		c.RawBytes, _ = ioutil.ReadAll(gzr)
	}
	c.Raw = ""
	c.Compressed = ""
	l.content[c.Path] = c
	for d := path.Dir(c.Path); d != "/"; d = path.Dir(d) {
		if _, ok := l.dirs[d]; !ok {
			l.dirs[d] = embeddedDir{name: d}
		}
	}
}

func (l *EmbeddedLoader) GetContents(path string) ([]byte, error) {
	if c, ok := l.content[path]; ok {
		return c.RawBytes, nil
	}
	return nil, os.ErrNotExist
}
