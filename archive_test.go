package archive

import (
	"archive/tar"
	"bytes"
	"io"

	"testing"
)

func check(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func assert(t *testing.T, cond bool, msg string) {
	if !cond {
		t.Fatal(msg)
	}
}

func TestTar(t *testing.T) {

	buff := &bytes.Buffer{}

	a := NewTarWriter(buff)
	r := tar.NewReader(buff)

	check(t, a.AddBytes("foo.txt", []byte("hello foo!")))
	check(t, a.AddBytes("bar.txt", []byte("hello bar!")))
	check(t, a.Close())

	//check for valid tar archive

	h, err := r.Next()
	check(t, err)
	assert(t, h.Name == "foo.txt", "first file should be foo.txt")
	b := make([]byte, h.Size)
	_, err = r.Read(b)
	check(t, err)
	assert(t, bytes.Compare(b, []byte("hello foo!")) == 0, "first file invalid")

	h, err = r.Next()
	check(t, err)
	assert(t, h.Name == "bar.txt", "second file should be bar.txt")
	b = make([]byte, h.Size)
	_, err = r.Read(b)
	check(t, err)
	assert(t, bytes.Compare(b, []byte("hello bar!")) == 0, "second file invalid")

	_, err = r.Next()
	assert(t, err == io.EOF, "should be end of file")
}
