package loader

import (
	"os"
	"path"
	"time"
)

type Content struct {
	Path            string
	Hash            string
	Modified        time.Time
	Raw             string
	Compressed      string
	RawBytes        []byte
	CompressedBytes []byte
}

func (f Content) Name() string {
	return path.Base(f.Path)
}

func (f Content) Size() int64 {
	return int64(len(f.Raw))
}

func (f Content) Mode() os.FileMode {
	return os.FileMode(0)
}

func (f Content) ModTime() time.Time {
	return f.Modified
}

func (f Content) IsDir() bool {
	return false
}

func (f Content) Sys() interface{} {
	return nil
}
