package main

import (
	"os"

	"github.com/jpillora/archive"
)

func main() {
	a := archive.NewZipWriter(os.Stdout)
	f, _ := os.Open("ping.txt")
	a.AddFile(f.Name(), f)
	f, _ = os.Open("pong.txt")
	a.AddFile(f.Name(), f)
	a.Close()
}
