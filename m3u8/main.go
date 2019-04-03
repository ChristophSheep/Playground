package main

import (
	"fmt"
	"time"

	"github.com/mysheep/cell"
)

const (
	DEBUG = true
)

func MediaGrapper(masterItems <-chan DownloadItem, mediaItems chan<- DownloadItem) {
	for {
		item := <-masterItems
		baseUrl := getBaseUrl(item.url)

		mediaUrl, err := getMediaPlayListUrl(item.url)
		if err == nil {
			if isRelativeUrl(mediaUrl) {
				mediaUrl = baseUrl + mediaUrl
			}
			mediaItems <- DownloadItem{url: mediaUrl, folder: item.folder}
		}
	}
}

func SegmentsDownloader(itemsQueue <-chan DownloadItem, downloaded chan<- m3u8URL) {

	downloadedUrls := map[m3u8URL]bool{}

	for {
		item := <-itemsQueue

		_, present := downloadedUrls[item.url]
		if present {
			continue
		}
		downloadedUrls[item.url] = false // queued but download not finished

		if DEBUG {
			// TODO: print to channel of console (stdout)

			//func Println(a ...interface{}) (n int, err error) {
			//	return Fprintln(os.Stdout, a...)
			//}

			//var (
			//	Stdin  = NewFile(uintptr(syscall.Stdin), "/dev/stdin")
			//	Stdout = NewFile(uintptr(syscall.Stdout), "/dev/stdout")
			//	Stderr = NewFile(uintptr(syscall.Stderr), "/dev/stderr")
			//)

			fmt.Println("Download url:", item.url)
			time.Sleep(5 * time.Second)

		} else {
			downloadItem(item)
		}
		downloadedUrls[item.url] = true // queued and finished
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
	segItemsQueue := make(chan DownloadItem, 100)
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
	go SegmentsGrapper(mediaItems, segItemsQueue)
	go SegmentsDownloader(segItemsQueue, downloaded) // TODO: more worker e.g. 3

	// Wait until quit
	//
	<-quit
}
