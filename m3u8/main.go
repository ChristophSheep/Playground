package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/grafov/m3u8"
)

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

	case m3u8.MEDIA:
		mediapl := p.(*m3u8.MediaPlaylist)
		fmt.Println("m3u8.MEDIA:")
		for i := uint(0); i < mediapl.Count(); i++ {
			s := mediapl.Segments[i]
			fmt.Println("i:", i, "uri:", s.URI)
		}
	}
}

func getFilename(urlRaw string) string {
	url, err := url.Parse(urlRaw)
	if err != nil {
		panic(err)
	}
	return path.Base(url.Path)
}

//
//   -->[ m3u8 ]-->
//

func main() {

	urlRaw := "http://orf1.orfstg.cdn.ors.at/out/u/orf1/q6a/manifest.m3u8"
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

	// given:
	//  - url (orf1 livestream)
	//  - timeStart
	//  - timeEnd

	// First load the mediaplaylist from masterplaylist
	// Always load mediaplaylist

	/*
		sequenceStartNr := 200756
		templateUrl := "http://orf1.orfstg.cdn.ors.at/out/u/orf1/q6a/manifest_4_###.ts?m=1552488594&f=5"

		i := sequenceStartNr
		for {
			url, tsFileName := getCurl(i, templateUrl)
			fmt.Println("url:", url, "filename:", tsFileName)
			i++
		}
	*/

}
