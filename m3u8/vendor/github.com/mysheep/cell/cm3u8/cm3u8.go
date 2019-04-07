package cm3u8

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"

	"github.com/grafov/m3u8"
)

// Master
//  - Variant 0 quality, url, ...
//  - Variant 1
//  - Variant 2

// Media
//  - Segment 0 .ts
//  - Segment 1 .ts
//  - Segment 2 .ts

type M3U8URL string

func MasterLoader(urls <-chan M3U8URL, masterPlayLists chan<- m3u8.MasterPlaylist) {

	for {

		m3u8Url := <-urls

		pl, listType, err := getPlaylist(m3u8Url)
		if err != nil {
			panic(err)
		}

		if listType == m3u8.MASTER {
			masterpl := pl.(*m3u8.MasterPlaylist)
			masterPlayLists <- *masterpl
		}
	}
}

func MediaLoader(urls <-chan M3U8URL, mediaPlayLists chan<- m3u8.MediaPlaylist) {
	m3u8Url := <-urls

	pl, listType, err := getPlaylist(m3u8Url)
	if err != nil {
		panic(err)
	}

	if listType == m3u8.MEDIA {
		mediapl := pl.(*m3u8.MediaPlaylist)
		mediaPlayLists <- *mediapl
	}

}

// getPlaylist get playlist from url master or media playlist
func getPlaylist(m3u8Url M3U8URL) (m3u8.Playlist, m3u8.ListType, error) {
	resp, err := http.Get(string(m3u8Url))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	reader := bytes.NewReader(body)
	return m3u8.DecodeFrom(reader, true)
}

// isRelativeUrl checks if the url is
//  - relativ or
//  - absolute
// url.
//
// e.g. http://server.com/file.ts is absolute url
//      file.ts is a relative url
func isRelativeUrl(urlRaw M3U8URL) bool {

	url, err := url.Parse(string(urlRaw))
	if err != nil {
		panic(err)
	}

	return (url.IsAbs() == false)
}

// getBaseUrl gets the base from with filename
// e.g. from http://server.com/folder/file.txt
//      get baseUrl -> http://server.com/folder/
func getBaseUrl(urlRaw M3U8URL) M3U8URL {

	url, err := url.Parse(string(urlRaw))
	if err != nil {
		panic(err)
	}

	//	[scheme:][//[userinfo@]host][/]path[?query][#fragment]

	res := url.Scheme + "://" + url.Host + path.Dir(url.Path) + "/"
	return M3U8URL(res)
}
