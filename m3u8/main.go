package main

import (
	"fmt"
	"time"

	"github.com/grafov/m3u8"
	"github.com/mysheep/cell"
	"github.com/mysheep/cell/cm3u8"
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
				printMsg("Grapper", fmt.Sprintf("Already in list - urlToDownload: %s ", urlToDownload))
			}
		}
	}

	stop := false

	// Wait for media playlist with segments
	// Download segments until stop signal is received
	//
	go func() {
		for {

			mediaUrl := makeAbsolute(baseUrl, cm3u8.M3U8URL(mediaUrlStr))
			mediaUrls <- mediaUrl

			mediaPlaylist := <-mediaPlaylists
			downloadSegments(mediaPlaylist)

			time.Sleep(time.Duration(mediaPlaylist.TargetDuration) * time.Second)

			//
			// TODO: Find out if last 3 segments are all downloaded
			//
			if stop {
				break
			}

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

	// Wait for Stop Signal
	//
	go func() {
		<-stopSignal
		printMsg("Grapper", "STOP signal was fired. Set stop = true.")
		stop = true
	}()
}

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
				start: time.Now().Add(2 * time.Second),
				end:   time.Now().Add(32 * time.Second),
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
