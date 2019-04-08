package main

import (
	"fmt"
	"time"

	"github.com/grafov/m3u8"
	"github.com/mysheep/cell"
	"github.com/mysheep/cell/cm3u8"
	"github.com/mysheep/cell/ctime"
	"github.com/mysheep/cell/web"
)

const (
	DEBUG = false
)

// download what?
//  - channel ?
//  - when ?
//  - to folder ?

var (
	channels = map[string]cm3u8.M3U8URL{
		"orf1": cm3u8.M3U8URL("http://orf1.orfstg.cdn.ors.at/out/u/orf1/q6a/manifest.m3u8"),
		"orf2": cm3u8.M3U8URL("http://orf2.orfstg.cdn.ors.at/out/u/orf2/q4a/manifest.m3u8"),
	}
)

type TimeSlot struct {
	start time.Time
	end   time.Time
}

type DownloadOrder struct {
	channel  string
	timeSlot TimeSlot
	folder   string
}

type DownloadItem struct {
	url    cm3u8.M3U8URL
	folder string
}

func Grapper(orders <-chan DownloadOrder) {

	// List of queued items (entry is there but false)
	// and downloaded items (entry is there but true)
	downloadedItems := map[cm3u8.M3U8URL]bool{}

	// Setup Channels
	timeSlots := make(chan TimeSlot)
	startSignal := make(chan bool)
	stopSignal := make(chan bool)

	masterUrls := make(chan cm3u8.M3U8URL)
	mediaUrls := make(chan cm3u8.M3U8URL)
	masterPlaylists := make(chan m3u8.MasterPlaylist)
	mediaPlaylists := make(chan m3u8.MediaPlaylist)

	downloadItems := make(chan DownloadItem, 3)
	downloadedUrls := make(chan cm3u8.M3U8URL, 3)

	// Setup Network
	//
	go StartStopTimer(timeSlots, startSignal, stopSignal)

	go cm3u8.MasterLoader(masterUrls, masterPlaylists)
	go cm3u8.MediaLoader(mediaUrls, mediaPlaylists)

	go Downloader(downloadItems, downloadedUrls)

	// Get the order
	//
	order := <-orders

	// Set timer
	timeSlots <- order.timeSlot

	// Wait for Start Signal ...
	<-startSignal
	// ... then insert url to network

	masterUrl := channels[order.channel]
	baseUrl := cm3u8.GetBaseUrl(masterUrl)

	masterUrls <- masterUrl

	// Wait for master playlist
	masterPlaylist := <-masterPlaylists
	mediaUrlStr := masterPlaylist.Variants[0].URI

	downloadSegments := func(mediaPlaylist m3u8.MediaPlaylist) {
		for i := uint(0); i < mediaPlaylist.Count(); i++ { // TODO

			mediaSegment := mediaPlaylist.Segments[i]
			urlToDownload := makeAbsolute(baseUrl, cm3u8.M3U8URL(mediaSegment.URI))

			_, present := downloadedItems[urlToDownload]
			if present == false {
				downloadedItems[urlToDownload] = false
				downloadItems <- DownloadItem{
					url:    urlToDownload,
					folder: order.folder,
				}
			} else {
				printMsg("Grapper", fmt.Sprintf("urlToDownload: %s already in list", urlToDownload))
			}
		}
	}

	// Wait for media playlist with segments
	// Download segments until
	// stop signal is received
	//
	go func() {
		for {
			mediaUrl := makeAbsolute(baseUrl, cm3u8.M3U8URL(mediaUrlStr))
			printMsg("Grapper", fmt.Sprintf("send mediaUrl: %s", mediaUrl))
			mediaUrls <- mediaUrl
			mediaPlaylist := <-mediaPlaylists
			printMsg("Grapper", fmt.Sprintf("got playlist back for %s", mediaUrl))
			downloadSegments(mediaPlaylist)
		}
	}()

	// Downloaded urls are marked in map as downloaded
	//
	go func() {
		for {
			downloadedUrl := <-downloadedUrls
			printMsg("Grapper", fmt.Sprintf("downloadedUrl: %s", downloadedUrl))
			downloadedItems[downloadedUrl] = true
		}
	}()

	// Stop
	go func() {
		select {
		case <-stopSignal:
			printMsg("Grapper", "STOP signal was fired")
		}
	}()
}

func Downloader(items <-chan DownloadItem, downloaded chan<- cm3u8.M3U8URL) {

	urls := make(chan string)
	contents := make(chan []byte)
	filenames := make(chan string)
	bytess := make(chan []byte)
	savedFilenames := make(chan string)

	go web.Downloader(urls, contents)
	go web.Saver(filenames, bytess, savedFilenames)

	for {
		item := <-items

		// Send url and ...
		urls <- string(item.url)

		// ... wait for downloaded content
		content := <-contents

		// Send filename and bytess and ...
		fileName := item.folder + getFilename(item.url)
		filenames <- fileName
		bytess <- content

		// ... wait for file is saved
		<-savedFilenames
		downloaded <- item.url
	}
}

func StartStopTimer(timeSlots <-chan TimeSlot, startSignals chan<- bool, stopSignals chan<- bool) {

	starts := make(chan time.Time)
	ends := make(chan time.Time)

	go ctime.Timer(starts, startSignals)
	go ctime.Timer(ends, stopSignals)

	for {
		ts := <-timeSlots
		starts <- ts.start
		ends <- ts.end
	}

}

/*
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
*/

func main() {

	dateFormat := "2006-01-02 15:04"
	startDateStr := "2019-04-06 15:04"
	startDateTime, err := time.Parse(dateFormat, startDateStr)
	if err == nil {
		fmt.Println(startDateTime)
	}

	//
	// Channels
	//
	quit := make(chan bool)
	orders := make(chan DownloadOrder)

	//
	// Commands function of console
	//
	dl := func() {

		// Create an order
		//
		dlo := DownloadOrder{
			channel: "orf2",
			timeSlot: TimeSlot{
				start: time.Now().Add(3 * time.Second),
				end:   time.Now().Add(20 * time.Second),
			},
			folder: "./download/",
		}

		// Send order into network
		//
		orders <- dlo

	}
	cmds := map[string]func(){
		"quit": func() { quit <- true },
		"dl":   func() { dl() },
	}
	go cell.Console(cmds) // stdout, stdin
	go Grapper(orders)

	// Wait until quit
	//
	<-quit
	printMsg("Application", "Quit now!!")
}
