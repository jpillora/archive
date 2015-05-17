package main

import (
	"os"

	"github.com/jpillora/archive"
)

func main() {
	a := archive.NewTarWriter(os.Stdout)
	a.AddBytes("foo.txt", []byte("hello foo!"))
	a.AddBytes("dir/bar.txt", []byte("hello bar!"))
	a.Close()
}
