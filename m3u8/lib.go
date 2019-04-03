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
)

type m3u8URL string

type DownloadItem struct {
	url    m3u8URL
	folder string
}

func getFilename(urlRaw m3u8URL) string {

	url, err := url.Parse(string(urlRaw))
	if err != nil {
		panic(err)
	}

	return path.Base(url.Path)
}

// isRelativeUrl checks if the url is
//  - relativ or
//  - absolute
// url.
//
// e.g. http://server.com/file.ts is absolute url
//      file.ts is a relative url
func isRelativeUrl(urlRaw m3u8URL) bool {

	url, err := url.Parse(string(urlRaw))
	if err != nil {
		panic(err)
	}

	return (url.IsAbs() == false)
}

// getBaseUrl gets the base from with filename
// e.g. from http://server.com/folder/file.txt
//      get baseUrl -> http://server.com/folder/
func getBaseUrl(urlRaw m3u8URL) m3u8URL {

	url, err := url.Parse(string(urlRaw))
	if err != nil {
		panic(err)
	}

	//	[scheme:][//[userinfo@]host][/]path[?query][#fragment]

	res := url.Scheme + "://" + url.Host + path.Dir(url.Path) + "/"
	return m3u8URL(res)
}

// downloadTo downloads the segments from a media file
// e.g. foo1001.ts
func downloadItem(item DownloadItem) {

	fmt.Println("download", item.url, "to folder", item.folder)

	resp, err := http.Get(string(item.url))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	filename := getFilename(item.url)
	fullpath := path.Join(item.folder, filename)

	f, err := os.Create(fullpath)
	if err != nil {
		panic(err)
	}
	defer f.Close() // f.Close will run when we're finished.
	f.Write(body)
}

// getPlaylist get playlist from url master or media playlist
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

// getMediaPlayListUrl get the media play list url
// from first variant of master playlist
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

// getMediaSegmentsUrls get the urls of the
// segments in a media play list file
func getMediaSegmentsUrls(m3u8Url m3u8URL) (urls []m3u8URL, err error) {

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
