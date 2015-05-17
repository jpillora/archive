package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jpillora/archive"
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
	//write .zip archive directly into response
	a := archive.NewZipWriter(w)
	f, _ := os.Open(LARGE_FILE)
	log.Println("sending...")
	t0 := time.Now()
	a.AddFile("file.bin", f)
	a.Close()
	log.Printf("sent in %s", time.Now().Sub(t0))

}
