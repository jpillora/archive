package main

import (
	"io"
	"os"

	"github.com/jpillora/archiver"
)

func main() {
	a := archiver.NewTarGz() //or NewTar() or NewZip()
	a.AddBytes("foo.txt", []byte("hello foo!"))
	a.AddBytes("dir/bar.txt", []byte("hello bar!"))
	a.Close()
	io.Copy(os.Stdout, a)
}
