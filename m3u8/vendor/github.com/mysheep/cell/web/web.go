package web

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type HttpGetter interface {
	httpGet(string) (*http.Response, error)
}

type MyGetter struct{}

func (g MyGetter) httpGet(url string) (resp *http.Response, err error) {
	return nil, nil
}

type FakeGetter struct{}

// see /usr/local/go/src/net/http/http.go line 100

// NoBody is an io.ReadCloser with no bytes. Read always returns EOF
// and Close always returns nil. It can be used in an outgoing client
// request to explicitly signal that a request has zero bytes.
// An alternative, however, is to simply set Request.Body to nil.
/*
var NoBody = noBody{}

type noBody struct{}

func (noBody) Read([]byte) (int, error)         { return 0, io.EOF }
func (noBody) Close() error                     { return nil }
func (noBody) WriteTo(io.Writer) (int64, error) { return 0, nil }

var (
	// verify that an io.Copy from NoBody won't require a buffer:
	_ io.WriterTo   = NoBody
	_ io.ReadCloser = NoBody
)

*/

type fakeReadCloser struct {
	url string
}

func (fakeReadCloser) Read(p []byte) (int, error) {
	// Read bytes into p
	return 0, io.EOF
}
func (fakeReadCloser) Close() error {
	return nil
}

func (g FakeGetter) httpGet(url string) (resp *http.Response, err error) {
	r := http.Response{
		Status:           "200 OK",
		StatusCode:       200,
		Proto:            "HTTP/1.0",
		ProtoMajor:       1,
		ProtoMinor:       0,
		Header:           http.Header{},
		Body:             fakeReadCloser{},
		ContentLength:    7, // TODO
		TransferEncoding: []string{},
		Close:            true,
		Uncompressed:     true,
		Trailer:          http.Header{},
		Request:          &http.Request{},
		TLS:              nil,
	}
	return &r, nil
}

// Downloader download content of given url and
// send it to output contents channel
func Downloader(urls <-chan string, contents chan<- []byte) {

	getBytes := func(url string) []byte {

		resp, err := http.Get(url) // HOW TO MOOK THAT ???
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		return bytes
	}

	for {
		url := <-urls
		bytes := getBytes(url)
		contents <- bytes
	}
}

// Saver save the bytes by given filename and
// send the save file name to output channel
func Saver(filenames <-chan string, bytess <-chan []byte, savedFilenames chan<- string) {

	createAndWrite := func(filename string, bytes []byte) {
		file, err := os.Create(filename)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		file.Write(bytes)
	}

	for {
		filename := <-filenames
		bytes := <-bytess
		createAndWrite(filename, bytes)
		savedFilenames <- filename
	}
}
