package ioutils

import (
	"bytes"
	"io/ioutil"
	"os"
)

// NewReaderFromFile return an io.Reader from file
func NewReaderFromFile(file string) *bytes.Reader {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	return bytes.NewReader(b)
}
