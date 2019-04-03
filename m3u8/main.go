package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/grafov/m3u8"
	"github.com/mysheep/cell"
)

type m3u8URL string

type DownloadItem struct {
	url    m3u8URL
	folder string
}

func downloadTo(url m3u8URL, folder string) {

	fmt.Println("download", url, "to", folder)

	resp, err := http.Get(string(url))
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

func Downloader(items <-chan DownloadItem, downloaded chan<- m3u8URL) {
	for {
		item := <-items
		// TODO: download
		fmt.Println("Download url:", item.url)
		//
		downloaded <- item.url
	}
}

func getPlaylist(m3u8Url m3u8URL) (m3u8.Playlist, m3u8.ListType, error) {
	resp, err := http.Get(string(m3u8Url))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	reader := bytes.NewReader(body)
	return m3u8.DecodeFrom(reader, true)
}

func getMediaPlayListUrl(m3u8Url m3u8URL) (uri m3u8URL, err error) {

	pl, listType, err := getPlaylist(m3u8Url)
	if err != nil {
		panic(err)
	}

	if listType == m3u8.MASTER {
		masterpl := pl.(*m3u8.MasterPlaylist)
		url := masterpl.Variants[0].URI
		return m3u8URL(url), nil
	}

	return m3u8URL(""), errors.New("m3u8 file is not a playlist of type MASTER")
}

func getMediaPlayListSegementsUrls(m3u8Url m3u8URL) (urls []m3u8URL, err error) {

	mapUrl := func(count uint, ss []*m3u8.MediaSegment, f func(m3u8.MediaSegment) m3u8URL) []m3u8URL {
		urls := make([]m3u8URL, count)
		var i uint = 0
		for i = 0; i < count; i++ {
			urls[i] = f(*ss[i])
		}
		return urls
	}

	getURL := func(segment m3u8.MediaSegment) m3u8URL {
		return m3u8URL(segment.URI)
	}

	pl, listType, err := getPlaylist(m3u8Url)
	if err != nil {
		panic(err)
	}

	if listType == m3u8.MEDIA {
		mediapl := pl.(*m3u8.MediaPlaylist)
		c := mediapl.Count()
		urls := mapUrl(c, mediapl.Segments, getURL)
		return urls, nil
	}

	empty := make([]m3u8URL, 0)
	return empty, errors.New("m3u8 is not a playlist file of type MEDIA")
}

func getFilename(urlRaw m3u8URL) string {

	url, err := url.Parse(string(urlRaw))
	if err != nil {
		panic(err)
	}

	return path.Base(url.Path)
}

func isRelativeUrl(urlRaw m3u8URL) bool {

	url, err := url.Parse(string(urlRaw))
	if err != nil {
		panic(err)
	}

	return (url.IsAbs() == false)
}

func getBaseUrl(urlRaw m3u8URL) m3u8URL {

	url, err := url.Parse(string(urlRaw))
	if err != nil {
		panic(err)
	}

	//	[scheme:][//[userinfo@]host][/]path[?query][#fragment]

	res := url.Scheme + "//" + url.Host + path.Dir(url.Path) + "/"
	return m3u8URL(res)
}

func main() {

	quit := make(chan bool)
	items := make(chan DownloadItem, 100)
	downloaded := make(chan m3u8URL, 100)

	urlMasterRaw := m3u8URL("http://orf1.orfstg.cdn.ors.at/out/u/orf1/q6a/manifest.m3u8")

	addItem := func(url m3u8URL, baseUrl m3u8URL) {

		if isRelativeUrl(url) {
			url = baseUrl + url
		}

		items <- DownloadItem{url: url, folder: "."}
	}

	test := func(urlRaw m3u8URL) {

		baseUrl := getBaseUrl(urlMasterRaw)

		url, err1 := getMediaPlayListUrl(m3u8URL(urlMasterRaw))
		if err1 == nil {
			urls, err2 := getMediaPlayListSegementsUrls(url)
			if err2 == nil {
				for _, url := range urls {
					addItem(url, baseUrl)
				}

			}
		}
	}

	//
	// Console Commands
	//
	cmds := map[string]func(){
		"quit": func() { quit <- true },
		"test": func() { test(urlMasterRaw) },
	}

	go Downloader(items, downloaded)
	go cell.Console(cmds)

	<-quit

}

/*
func ReadM3U8() {

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
*/
