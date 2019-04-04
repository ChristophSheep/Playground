package main

import (
	"fmt"
	"time"

	"github.com/mysheep/cell"
	"github.com/mysheep/cell/boolean"
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
		time.Sleep(1 * time.Second)
	}
}

func SegmentsGrapper(mediaItems <-chan DownloadItem, quitRequest <-chan bool, segmentItems chan<- DownloadItem) {

	quit := false
	go func() {
		for {
			<-quitRequest
			quit = true
		}
	}()

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
		fmt.Println("Grap", len(segUrls), " segment(s). Queue them.")

		if err == nil {
			for _, surl := range segUrls {
				queueItem(surl)
			}
		}

		if quit {
			fmt.Println("Grabber - quit requested")
			break
		}

		time.Sleep(1 * time.Second)
	}
}

func SegmentsDownloader(itemsQueue <-chan DownloadItem, quitRequest <-chan bool, downloaded chan<- m3u8URL) {

	quit := false

	go func() {
		for {
			<-quitRequest // wait for quit request
			quit = true
		}
	}()

	downloadedUrls := map[m3u8URL]bool{}

	for {
		item := <-itemsQueue

		_, present := downloadedUrls[item.url]
		if present {
			continue
		}
		downloadedUrls[item.url] = false // queued but download not finished

		if DEBUG {
			fmt.Println("Download url:", item.url)
			time.Sleep(5 * time.Second)
		} else {
			downloadItem(item)
		}
		downloadedUrls[item.url] = true // queued and finished
		downloaded <- item.url

		if quit {
			fmt.Println("Downloader - quit requested")
			break
		}
	}
}

func main() {

	//
	// Channels
	//
	quit := make(chan bool)
	quitRequest0 := make(chan bool)
	quitRequest1 := make(chan bool)
	quitRequest2 := make(chan bool)

	masterItems := make(chan DownloadItem)
	mediaItems := make(chan DownloadItem)
	segItemsQueue := make(chan DownloadItem, 100)
	downloaded := make(chan m3u8URL, 100)

	//
	// Commands function of console
	//
	// urlMasterPlaylistOrf1 := m3u8URL("http://orf1.orfstg.cdn.ors.at/out/u/orf1/q6a/manifest.m3u8")
	urlMasterPlaylistOrf2 := m3u8URL("http://orf2.orfstg.cdn.ors.at/out/u/orf2/q4a/manifest.m3u8")
	dl := func() {
		masterItems <- DownloadItem{url: urlMasterPlaylistOrf2, folder: "."}
	}
	cmds := map[string]func(){
		"quit": func() { quit <- true },
		"qr":   func() { quitRequest0 <- true },
		"dl":   func() { dl() },
	}
	go cell.Console(cmds) // stdout, stdin

	//
	// Setup network
	//
	go boolean.Distributor(quitRequest0, quitRequest1, quitRequest2)
	go MediaGrapper(masterItems, mediaItems)
	go SegmentsGrapper(mediaItems, quitRequest1, segItemsQueue)
	go SegmentsDownloader(segItemsQueue, quitRequest2, downloaded) // TODO: more worker e.g. 3

	// Wait until quit
	//
	<-quit
}

// TODO: print to channel of console (stdout)

//func Println(a ...interface{}) (n int, err error) {
//	return Fprintln(os.Stdout, a...)
//}

//var (
//	Stdin  = NewFile(uintptr(syscall.Stdin), "/dev/stdin")
//	Stdout = NewFile(uintptr(syscall.Stdout), "/dev/stdout")
//	Stderr = NewFile(uintptr(syscall.Stderr), "/dev/stderr")
//)
