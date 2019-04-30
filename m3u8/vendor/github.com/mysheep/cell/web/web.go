package web

import (
	"io/ioutil"
	"net/http"
	"os"
)

// Downloader download content of given url and
// send it to output contents channel
func Downloader(urls <-chan string, contents chan<- []byte) {

	getBytes := func(url string) []byte {

		resp, err := http.Get(url)
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
