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

func SegmentsGrapper(mediaItems <-chan DownloadItem, startSignal, stopSignal <-chan bool, segmentItems chan<- DownloadItem) {

	<-startSignal // Wait for start

	quit := false
	go func() {
		for {
			<-stopSignal
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

		segUrls, err, targetDurationInSec := getMediaSegmentsUrls(item.url)
		msg := fmt.Sprintf("Grap %v segment(s) and queue them to download", len(segUrls))
		printMsg("Grabber", msg)

		if err == nil {
			for _, surl := range segUrls {
				queueItem(surl)
			}
		}

		if quit {
			printMsg("Grabber", "Quit requested, stop grabbing")
			break
		}

		// TODO: HOW LONG TO WAIT ????
		time.Sleep(time.Duration(targetDurationInSec/2.0) * time.Second)
	}
}

func SegmentsDownloader(itemsQueue <-chan DownloadItem, downloaded chan<- m3u8URL) {

	// TODO: Rethink - Give other component responsibility
	downloadedUrls := map[m3u8URL]bool{}

	// TODO: Rethink - Give other component responsibility
	isAlreadyDownloaded := func(url m3u8URL) bool {
		isDownloaded, present := downloadedUrls[url]
		return present && isDownloaded
	}

	for {

		item := <-itemsQueue

		if isAlreadyDownloaded(item.url) {
			continue
		}

		printMsg("Downloader", fmt.Sprintf("Start    download '%s'", item.url))
		downloadItem(item)
		printMsg("Downloader", fmt.Sprintf("Finished download '%s'", item.url))

		downloadedUrls[item.url] = true
		downloaded <- item.url

	}
}

func Timer(startTime time.Time, stopTime time.Time, startSignal chan<- bool, stopSignal chan<- bool) {

	// START
	//
	durationToStart := startTime.Sub(time.Now())
	printMsg("Timer", fmt.Sprintf("Wait to start at '%v' in '%v' seconds", startTime, durationToStart))
	time.Sleep(durationToStart)

	printMsg("Timer", fmt.Sprintf("Send start '%v'", startTime))
	startSignal <- true

	// STOP
	//
	durationToStop := stopTime.Sub(time.Now())
	time.Sleep(durationToStop)
	printMsg("Timer", fmt.Sprintf("Send stop '%v'", stopTime))

	stopSignal <- true

}

func Terminator(quitRequest <-chan bool, itemsQueue <-chan DownloadItem, quit chan<- bool) {

	shouldQuit := false

	go func() {
		for {
			<-quitRequest // wait for quit request
			shouldQuit = true
		}
	}()

	isEmpty := func() bool {
		return len(itemsQueue) == 0
	}

	for {
		if shouldQuit && isEmpty() {
			printMsg("Terminator", "Send quit signal")
			quit <- true
		}
		time.Sleep(10 * time.Millisecond)
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

	startSignal := make(chan bool)
	stopSignal := make(chan bool)

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
		start := time.Now().Add(5 * time.Second)
		stop := time.Now().Add(15 * time.Second)

		go Timer(start, stop, startSignal, stopSignal)

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

	go MediaGrapper(masterItems, mediaItems)
	go boolean.Distributor(stopSignal, quitRequest1, quitRequest2)
	go SegmentsGrapper(mediaItems, startSignal, quitRequest1, segItemsQueue)
	go SegmentsDownloader(segItemsQueue, downloaded) // TODO: more worker e.g. 3
	go Terminator(quitRequest2, segItemsQueue, quit)

	// Wait until quit
	//
	<-quit
	printMsg("Application", "Quit now!!")
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
