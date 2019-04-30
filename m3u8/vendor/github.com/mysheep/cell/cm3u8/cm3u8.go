package cm3u8

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/grafov/m3u8"
)

// Master
//  - Variant 0 {quality 0, media url 0, ...
//  - Variant 1 {quality 1, media url 1, ...)
//  - Variant 2 {quality 2, media url 2, ...)

// Media
//  - Segment 0.ts
//  - Segment 1.ts
//  - Segment 2.ts

type M3U8URL string

// MasterLoader waits on input channel to receive a url
// if a url is incoming he load the master playlist
// and send it to the outputs channel
func MasterLoader(urls <-chan M3U8URL, masterPlayLists chan<- m3u8.MasterPlaylist) {

	getMasterPlayList := func(url M3U8URL) (m3u8.MasterPlaylist, error) {
		empty := m3u8.MasterPlaylist{}
		pl, listType, err := getPlaylist(url)
		if err != nil {
			return empty, errors.New("Url could not be loaded")
		}
		if listType == m3u8.MASTER {
			masterpl := pl.(*m3u8.MasterPlaylist)
			return *masterpl, nil
		}
		return empty, errors.New("Url was not a masterplay m3u8")
	}

	for {
		url := <-urls
		masterPlaylist, err := getMasterPlayList(url)
		if err == nil {
			masterPlayLists <- masterPlaylist
		}
	}
}

// MediaLoader reveices urls of a media playlist
// on input channel, then loads the playlist and
// sends to the output channel
func MediaLoader(urls <-chan M3U8URL, mediaPlayLists chan<- m3u8.MediaPlaylist) {

	getMediaPlaylist := func(url M3U8URL) (m3u8.MediaPlaylist, error) {
		empty := m3u8.MediaPlaylist{}
		pl, listType, err := getPlaylist(url)
		if err != nil {
			return empty, errors.New("Url could not be loaded")
		}
		if listType == m3u8.MEDIA {
			mediapl := pl.(*m3u8.MediaPlaylist)
			return *mediapl, nil
		}
		return empty, errors.New("Url is not a mediaplaylist")
	}

	for {
		url := <-urls
		mediaPlayList, err := getMediaPlaylist(url)
		if err == nil {
			mediaPlayLists <- mediaPlayList
		}
	}
}

// Repeater wait to receive an url on input channel
// if he has an url received he repeat to send
// the url to output channel with the given interval in seconds
func Repeater(intervalInSec uint, ins <-chan M3U8URL, outs chan<- M3U8URL) {

	val := M3U8URL("")
	valIsSet := false

	go func() {
		for {
			val = <-ins
			valIsSet = true
		}
	}()

	for {
		time.Sleep(time.Duration(intervalInSec) * time.Second)
		//fmt.Println("Repeater", "-", "interval:", interval)
		if valIsSet {
			//fmt.Println("Repeater", "-", "send val:", val)
			outs <- val
		}
	}
}

// Switch receive on or off signal in onOffs channel
// if the swithc receive a onOff = true
// the incoming urls and send to output channel
// else the incoming urls are consumed, but not send to output channel
func Switch(onOffs <-chan bool, ins <-chan M3U8URL, outs chan<- M3U8URL) {

	onOff := false
	go func() {
		for {
			onOff = <-onOffs // wait for signal
			//fmt.Println("Switch", "-", "onOff:", onOff)
		}
	}()

	for {
		val := <-ins
		//fmt.Println("Switch", "-", "onOff:", onOff)
		if onOff {
			outs <- val
			//fmt.Println("Switch", "-", "val:", val)
		}
	}

}

// SegmentsGrapper graps the url of each segment and send its to output
func SegmentsGrapper(mediaPlaylists <-chan m3u8.MediaPlaylist, mediaSegmentURIs chan<- M3U8URL) {
	for {
		mediaPlaylist := <-mediaPlaylists
		//fmt.Println("SegmentsGrapper", "-", "mediaPlaylist count:", mediaPlaylist.Count())
		for i := uint(0); i < mediaPlaylist.Count(); i++ {
			mediaSegment := mediaPlaylist.Segments[i]
			//fmt.Println("SegmentsGrapper", "-", "mediaPlaylist segment i:", i, " uri:", mediaSegment.URI)
			mediaSegmentURIs <- M3U8URL(mediaSegment.URI)
		}
	}
}

// Mapper maps incoming urls by calling res = fn(url) and send the res to outs
func Mapper(ins <-chan M3U8URL, outs chan<- M3U8URL, fn func(M3U8URL) M3U8URL) {
	for {
		val := <-ins
		outs <- fn(val)
	}
}

// Filter incoming url, if fn(url) is true is passes the filter else not
func Filter(ins <-chan M3U8URL, outs chan<- M3U8URL, fn func(M3U8URL) bool) {
	for {
		val := <-ins
		if fn(val) {
			//fmt.Println("Filter", "-", "url send out:", val)
			outs <- val
		}
	}
}

// getPlaylist get playlist from url master or media playlist
func getPlaylist(m3u8Url M3U8URL) (m3u8.Playlist, m3u8.ListType, error) {
	resp, err := http.Get(string(m3u8Url))
	if err != nil {
		fmt.Println(err)
		return nil, m3u8.MASTER, err
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
func IsRelativeUrl(urlRaw M3U8URL) bool {

	url, err := url.Parse(string(urlRaw))
	if err != nil {
		panic(err)
	}

	return (url.IsAbs() == false)
}

// getBaseUrl gets the base from with filename
// e.g. from http://server.com/folder/file.txt
//      get baseUrl -> http://server.com/folder/
func GetBaseUrl(urlRaw M3U8URL) M3U8URL {

	url, err := url.Parse(string(urlRaw))
	if err != nil {
		panic(err)
	}

	//	[scheme:][//[userinfo@]host][/]path[?query][#fragment]

	res := url.Scheme + "://" + url.Host + path.Dir(url.Path) + "/"
	return M3U8URL(res)
}
