package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/jpillora/archiver"
)

func main() {
	a := archiver.NewTar()

	//start POSTing straight away
	wait := make(chan bool)
	go func() {
		resp, err := http.Post("https://echo.jpillora.com/archiver-stream", "application/x-gzip", a)
		if err != nil {
			log.Fatal("post failed")
		}
		b, _ := ioutil.ReadAll(resp.Body)
		log.Printf("response: %s\n%s", resp.Status, b)
		close(wait)
	}()

	//simulate "searching" for files
	delay()
	a.AddBytes("foo.txt", []byte("hello foo!"))
	delay()
	a.AddBytes("dir/bar.txt", []byte("hello bar!"))
	delay()
	a.AddBytes("dir/zip/bazz.txt", []byte("hello bazz!"))
	delay()

	//close will finalize the archive and send EOF to http.Post so it knows the request is finished
	a.Close()
	<-wait
}

func delay() {
	time.Sleep(300 * time.Millisecond)
}
