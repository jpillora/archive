package archiver

import (
	"compress/gzip"
	"errors"
	"io"

	"os"
	"path/filepath"
)

const DefaultMaxDirSize = 4e9  //4GB
const DefaultMaxDirFiles = 1e5 //100,000

//Archiver is a higher level API over archive/zip,tar
type Archiver struct {
	//config
	DirMaxSize  int64
	DirMaxFiles int
	//state
	dst     io.Writer
	archive archive
}

//NewWriter is useful when you have a dynamic archive filename with extension
func NewWriter(filename string, dst io.Writer) (*Archiver, error) {

	a := &Archiver{
		DirMaxSize:  DefaultMaxDirSize,
		DirMaxFiles: DefaultMaxDirFiles,
		dst:         dst,
	}

	switch Extension(filename) {
	case ".tar":
		a.archive = newTarArchive(dst)
	case ".tar.gz":
		gz := gzip.NewWriter(dst)
		a.dst = gz
		a.archive = newTarArchive(gz)
	case ".zip":
		a.archive = newZipArchive(dst)
	default:
		return nil, errors.New("Invalid extension: " + filename)
	}

	return a, nil

}

func NewTarWriter(dst io.Writer) *Archiver {
	a, _ := NewWriter(".tar", dst)
	return a
}

func NewTarGzWriter(dst io.Writer) *Archiver {
	a, _ := NewWriter(".tar.gz", dst)
	return a
}

func NewZipWriter(dst io.Writer) *Archiver {
	a, _ := NewWriter(".zip", dst)
	return a
}

func (a *Archiver) AddBytes(path string, contents []byte) error {
	if err := checkPath(path); err != nil {
		return err
	}
	return a.archive.addBytes(path, contents)
}

func (a *Archiver) AddFile(path string, f *os.File) error {
	info, err := f.Stat()
	if err != nil {
		return err
	}
	return a.AddInfoFile(path, info, f)
}

//You can prevent archiver from performing an extra Stat by using AddInfoFile
//instead of AddFile
func (a *Archiver) AddInfoFile(path string, info os.FileInfo, f *os.File) error {
	if err := checkPath(path); err != nil {
		return err
	}
	return a.archive.addFile(path, info, f)
}

func (a *Archiver) AddDir(path string) error {
	size := a.DirMaxSize
	num := a.DirMaxFiles
	return filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if !info.Mode().IsRegular() {
			return nil
		}
		size -= info.Size()
		if size <= 0 {
			return errors.New("Surpassed maximum archive size")
		}
		num--
		if num <= 0 {
			return errors.New("Surpassed maximum number of files in archive")
		}
		// log.Printf("#%d %09d add file %s", num, size, p)
		rel, err := filepath.Rel(path, p)
		if err != nil {
			return err
		}
		f, err := os.Open(p)
		if err != nil {
			return err
		}
		defer f.Close()
		return a.archive.addFile(rel, info, f)
	})
}

func (a *Archiver) Close() error {
	if err := a.archive.close(); err != nil {
		return err
	}
	if gz, ok := a.dst.(*gzip.Writer); ok {
		if err := gz.Close(); err != nil {
			return err
		}
	}
	return nil
}
