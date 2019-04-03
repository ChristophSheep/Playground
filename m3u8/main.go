package main

import (
	"fmt"

	"github.com/mysheep/cell"
)

const (
	DEBUG = true
)

func MediaGrapper(masterItems <-chan DownloadItem, mediaItems chan<- DownloadItem) {
	for {
		item := <-masterItems
		mediaUrl, err := getMediaPlayListUrl(item.url)
		if err == nil {
			mediaItems <- DownloadItem{url: mediaUrl, folder: item.folder}
		}
	}
}

func SegmentsDownloader(items <-chan DownloadItem, downloaded chan<- m3u8URL) {
	for {
		item := <-items
		if DEBUG {
			fmt.Println("Download url:", item.url)
		} else {
			downloadItem(item)
		}
		downloaded <- item.url
	}
}

func SegmentsGrapper(mediaItems <-chan DownloadItem, segmentItems chan<- DownloadItem) {

	item := <-mediaItems
	baseUrl := getBaseUrl(item.url)

	queueItem := func(surl m3u8URL) {
		if isRelativeUrl(surl) {
			surl = baseUrl + surl
		}
		segmentItems <- DownloadItem{url: surl, folder: item.folder}
	}

	for {
		segUrls, err := getMediaSegmentsUrls(item.url)
		if err == nil {
			for _, surl := range segUrls {
				//
				// TODO: Check if segments already queued and/or downloaded
				//
				queueItem(surl)
			}
		}
	}
}

func main() {

	//
	// Channels
	//
	quit := make(chan bool)

	masterItems := make(chan DownloadItem)
	mediaItems := make(chan DownloadItem)
	segItems := make(chan DownloadItem, 3)
	downloaded := make(chan m3u8URL, 100)

	//
	// Commands function of console
	//
	urlMasterPlaylistOrf1 := m3u8URL("http://orf1.orfstg.cdn.ors.at/out/u/orf1/q6a/manifest.m3u8")
	dl := func() {
		masterItems <- DownloadItem{url: urlMasterPlaylistOrf1, folder: "."}
	}
	cmds := map[string]func(){
		"quit": func() { quit <- true },
		"dl":   func() { dl() },
	}
	go cell.Console(cmds) // stdout, stdin

	//
	// Setup network
	//
	go MediaGrapper(masterItems, mediaItems)
	go SegmentsGrapper(mediaItems, segItems)
	go SegmentsDownloader(segItems, downloaded) // TODO: more worker e.g. 3

	// Wait until quit
	//
	<-quit
}
