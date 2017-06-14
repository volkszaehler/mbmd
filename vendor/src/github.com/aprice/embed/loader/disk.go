package loader

import (
	"io/ioutil"
	"net/http"
)

type DiskLoader struct {
	Root string
	http.Handler
}

func NewOnDisk(root string) *DiskLoader {
	return &DiskLoader{
		Root:    root,
		Handler: http.FileServer(http.Dir(root)),
	}
}

func (l *DiskLoader) GetContents(path string) ([]byte, error) {
	return ioutil.ReadFile(l.Root + path)
}
