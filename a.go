package archive

import (
	"os"
	"time"
)

//a common interface to tar/zip
type archive interface {
	addBytes(path string, contents []byte, mtime time.Time) error
	addFile(path string, info os.FileInfo, f *os.File) error
	close() error
}
