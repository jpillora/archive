package archiver

import (
	"bytes"
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
	buff    *bytes.Buffer
	gz      *gzip.Writer
	archive archive
	closed  bool
}

//New is useful when you have the destination path and want to
//extract it's extension to use as the archive type
func New(path string) (*Archiver, error) {
	switch Extension(path) {
	case ".tar":
		return NewTar(), nil
	case ".tar.gz":
		return NewTarGz(), nil
	case ".zip":
		return NewZip(), nil
	}
	return nil, errors.New("Invalid extension: " + path)
}

func NewTar() *Archiver {
	buff := &bytes.Buffer{}
	return new(buff, newTarArchive(buff))
}

func NewTarGz() *Archiver {
	buff := &bytes.Buffer{}
	gz := gzip.NewWriter(buff)
	a := new(buff, newTarArchive(buff))
	a.gz = gz
	return a
}

func NewZip() *Archiver {
	buff := &bytes.Buffer{}
	return new(buff, newZipArchive(buff))
}

func new(buff *bytes.Buffer, a archive) *Archiver {
	return &Archiver{
		DirMaxSize:  DefaultMaxDirSize,
		DirMaxFiles: DefaultMaxDirFiles,
		buff:        buff,
		archive:     a,
	}
}

func (a *Archiver) Read(p []byte) (int, error) {
	n, err := a.buff.Read(p)
	if err == io.EOF && !a.closed {
		//TODO sleep here to slow loop?
		return n, nil
	}
	return n, err
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
	if a.gz != nil {
		if err := a.gz.Close(); err != nil {
			return err
		}
	}
	a.closed = true
	return nil
}
