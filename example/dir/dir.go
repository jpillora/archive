package main

import (
	"os"

	"github.com/jpillora/archiver"
)

func main() {
	a, _ := archiver.NewWriter("file.tar", os.Stdout) //detects .tar
	a.AddDir("foo")
	a.Close()
}
