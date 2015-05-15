package main

import (
	"io"
	"os"

	"github.com/jpillora/archiver"
)

func main() {
	a := archiver.NewZip()
	f, _ := os.Open("ping.txt")
	a.AddFile(f.Name(), f)
	f, _ = os.Open("pong.txt")
	a.AddFile(f.Name(), f)
	a.Close()
	io.Copy(os.Stdout, a)
}
