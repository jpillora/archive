package main

import (
	"io"
	"os"

	"github.com/jpillora/archiver"
)

func main() {
	a := archiver.NewTar()
	a.AddDir("foo")
	a.Close()
	io.Copy(os.Stdout, a)
}
