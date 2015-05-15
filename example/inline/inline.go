package main

import (
	"os"

	"github.com/jpillora/archiver"
)

func main() {
	a := archiver.NewTarWriter(os.Stdout)
	a.AddBytes("foo.txt", []byte("hello foo!"))
	a.AddBytes("dir/bar.txt", []byte("hello bar!"))
	a.Close()
}
