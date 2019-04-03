package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/grafov/m3u8"
)

type URL string

func downloadTo(url string, folder string) {

	fmt.Println("download", url, "to", folder)

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	filename := getFilename(url)
	fullpath := path.Join(folder, filename)

	f, err := os.Create(fullpath)
	if err != nil {
		panic(err)
	}
	defer f.Close() // f.Close will run when we're finished.
	f.Write(body)
}

func Downloader(nextUrl <-chan string, downloaded chan<- string) {
	for {
		url := <-nextUrl
		// TODO: download
		downloaded <- url
	}
}

func getMediaPlayListUrl(m3u8Url URL) (uri URL, err error) {

	resp, err := http.Get(string(m3u8Url))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	reader := bytes.NewReader(body)

	pl, listType, err := m3u8.DecodeFrom(reader, true)
	if err != nil {
		panic(err)
	}

	if listType == m3u8.MASTER {
		masterpl := pl.(*m3u8.MasterPlaylist)
		url := masterpl.Variants[0].URI
		return URL(url), nil
	}

	return URL(""), errors.New("m3u8 file is not a masterplaylist")
}

func getMediaPlayListSegments(m3u8Url URL) (urls []URL, err error) {

	mapUrl := func(c uint, ss []*m3u8.MediaSegment, f func(m3u8.MediaSegment) URL) []URL {
		fmt.Println("len:", c)
		vsm := make([]URL, c)
		var i uint = 0
		for i = 0; i < c; i++ {
			vsm[i] = f(*ss[i])
		}
		return vsm
	}

	resp, err := http.Get(string(m3u8Url))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	reader := bytes.NewReader(body)
	pl, listType, err := m3u8.DecodeFrom(reader, true)
	if err != nil {
		panic(err)
	}

	getURL := func(segment m3u8.MediaSegment) URL {
		return URL(segment.URI)
	}

	if listType == m3u8.MEDIA {
		mediapl := pl.(*m3u8.MediaPlaylist)
		c := mediapl.Count()
		urls := mapUrl(c, mediapl.Segments, getURL)
		return urls, nil
	}

	empty := make([]URL, 0)
	return empty, errors.New("m3u8 is not a media playlist file")
}

func getFilename(urlRaw string) string {
	url, err := url.Parse(urlRaw)
	if err != nil {
		panic(err)
	}
	return path.Base(url.Path)
}

func main() {

	urlRaw := "http://orf1.orfstg.cdn.ors.at/out/u/orf1/q6a/manifest.m3u8"
	url, err := getMediaPlayListUrl(URL(urlRaw))
	if err == nil {
		urls, err2 := getMediaPlayListSegments(url)
		if err2 != nil {
			panic(err)
		}
		for i, url := range urls {
			fmt.Println("i:", i, "url:", url)
			downloadTo(string(url), ".")
		}
	}

	/*
		url, err := url.Parse(urlRaw)
		if err != nil {
			panic(err)
		}
		fmt.Println("scheme:", url.Scheme)
		fmt.Println("host:", url.Host)
		fmt.Println("path:", url.Path)
		fmt.Println("isAbs:", url.IsAbs())
		fmt.Println("base path:", path.Base(url.Path))
		fmt.Println("dir path:", path.Dir(url.Path))
		fmt.Println("ext path:", path.Ext(path.Base(url.Path)))

		folder := "./"
		downloadTo(urlRaw, folder)

		fmt.Println()
		fmt.Println()
		ReadM3U8()

	*/

}

func ReadM3U8() {
	//(*File, error)
	f, err := os.Open("./example/masterplaylist.m3u8")
	if err != nil {
		panic(err)
	}
	p, listType, err := m3u8.DecodeFrom(bufio.NewReader(f), true)
	if err != nil {
		panic(err)
	}

	switch listType {

	case m3u8.MASTER:
		masterpl := p.(*m3u8.MasterPlaylist)
		fmt.Println("m3u8.MASTER:")
		fmt.Printf("%+v\n", masterpl)
		for i, v := range masterpl.Variants { // Variants of different streams
			fmt.Println("i:", i, "uri:", v.URI)
			fmt.Println("i:", i, "codecs:", v.Codecs)
			fmt.Println("i:", i, "resolution:", v.Resolution)
			fmt.Println("i:", i, "name:", v.Name)
		}

	case m3u8.MEDIA:
		mediapl := p.(*m3u8.MediaPlaylist) // Segments of a live stream
		fmt.Println("m3u8.MEDIA:")
		for i, s := range mediapl.Segments {
			fmt.Println("i:", i, "uri:", s.URI)
		}
	}
}
