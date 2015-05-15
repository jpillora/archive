package archiver

import "os"

//a common interface to tar/zip
type archive interface {
	addBytes(path string, contents []byte) error
	addFile(path string, info os.FileInfo, f *os.File) error
	close() error
}
