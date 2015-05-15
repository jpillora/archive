package main

import (
	"os"

	"github.com/jpillora/archiver"
)

func main() {
	a := archiver.NewZipWriter(os.Stdout)
	f, _ := os.Open("ping.txt")
	a.AddFile(f.Name(), f)
	f, _ = os.Open("pong.txt")
	a.AddFile(f.Name(), f)
	a.Close()
}
