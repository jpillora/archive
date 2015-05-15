package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jpillora/archiver"
)

const LARGE_FILE = "/Users/jpillora/file.mp4"

func main() {
	if _, err := os.Stat(LARGE_FILE); err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/", handle)
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/zip")
	w.WriteHeader(200)
	//write compression stream
	a := archiver.NewZipWriter(w)
	log.Println("sending...")
	f, _ := os.Open(LARGE_FILE)
	t0 := time.Now()
	a.AddFile("file.bin", f)
	a.Close()
	log.Printf("sent in %s", time.Now().Sub(t0))

}
