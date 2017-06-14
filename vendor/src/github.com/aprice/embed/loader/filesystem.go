package loader

import (
	"bytes"
	"net/http"
	"os"
	"path"
	"time"
)

// Open implements http.FileSystem
func (l *EmbeddedLoader) Open(name string) (http.File, error) {
	if v, ok := l.content[path.Clean(name)]; ok {
		return &embeddedFile{bytes.NewReader(v.RawBytes), v}, nil
	}
	if v, ok := l.dirs[path.Clean(name)]; ok {
		return v, nil
	}
	return nil, os.ErrNotExist
}

// embeddedFile implements http.File and is used once per request
type embeddedFile struct {
	*bytes.Reader
	*Content
}

func (f embeddedFile) Close() error {
	return nil
}

func (f embeddedFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

func (f embeddedFile) Stat() (os.FileInfo, error) {
	return f, nil
}

func (f embeddedFile) Size() int64 {
	return int64(len(f.RawBytes))
}

type embeddedDir struct {
	name    string
	modTime time.Time
}

func (f embeddedDir) Name() string {
	return path.Base(f.name)
}

func (f embeddedDir) Size() int64 {
	return 0
}

func (f embeddedDir) Mode() os.FileMode {
	return os.ModeDir
}

func (f embeddedDir) ModTime() time.Time {
	return f.modTime
}

func (f embeddedDir) IsDir() bool {
	return true
}

func (f embeddedDir) Sys() interface{} {
	return nil
}

func (f embeddedDir) Read(p []byte) (int, error) {
	return 0, os.ErrInvalid
}

func (f embeddedDir) Seek(offset int64, whence int) (int64, error) {
	return 0, os.ErrInvalid
}

func (f embeddedDir) Close() error {
	return nil
}

func (f embeddedDir) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

func (f embeddedDir) Stat() (os.FileInfo, error) {
	return f, nil
}
