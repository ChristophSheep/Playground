package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/grafov/m3u8"
	"github.com/mysheep/cell"
	"github.com/mysheep/cell/cm3u8"
)

const (
	DEBUG = false
)

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

func (ts TimeSlot) String() string {
	return fmt.Sprintf("{start:'%v', end:'%v'}", ts.start, ts.end)
}

type DownloadOrder struct {
	channel  string
	timeSlot TimeSlot
	folder   string
}

func (do DownloadOrder) String() string {
	return fmt.Sprintf("{channel:'%v', time:%v, folder:'%v'}", do.channel, do.timeSlot, do.folder)
}

type DownloadItem struct {
	url    cm3u8.M3U8URL
	folder string
}

func Grapper(orders <-chan DownloadOrder) {

	timeSlots := make(chan TimeSlot)
	startSignal := make(chan bool)
	stopSignal := make(chan bool)

	masterUrls := make(chan cm3u8.M3U8URL)
	masterPlaylists := make(chan m3u8.MasterPlaylist)
	mediaPlaylistUrls := make(chan cm3u8.M3U8URL)

	// Setup Network
	//
	go StartStopTimer(timeSlots, startSignal, stopSignal)
	go cm3u8.MasterLoader(masterUrls, masterPlaylists)
	go MediaLoaderSlidingWindow(mediaPlaylistUrls, startSignal, stopSignal)

	// Get the order
	//
	for {

		order := <-orders

		printMsg("Grapper", fmt.Sprint("order: ", order))
		masterUrl := channels[order.channel]
		baseUrl := cm3u8.GetBaseUrl(masterUrl)

		mediaPlaylistUrl := cm3u8.M3U8URL("")
		counter := 0
		maxTries := 10

		for {

			printMsg("Grapper", fmt.Sprintf("masterUrl: %s", masterUrl))

			// Set master url to download ..
			masterUrls <- masterUrl
			// .. and wait for masterPlaylist
			masterPlaylist := <-masterPlaylists

			mediaPlaylistUrl = cm3u8.M3U8URL(masterPlaylist.Variants[0].URI) // TODO CHeck
			mediaPlaylistUrl = makeAbsolute(baseUrl, mediaPlaylistUrl)

			// https://apasfiis.sf.apa.at/ipad/gp/livestream_Q6A.mp4/chunklist.m3u8?lbs=20190412132743573&origin=http%253a%252f%252fvarorfvod.sf.apa.at%252fsystem_clips%252flivestream_Q6A.mp4%252fchunklist.m3u8&ip=129.27.216.70&ua=Go-http-client%252f1.1

			if counter > 10 {
				break
			}

			if strings.Contains(string(mediaPlaylistUrl), "chunklist.m3u8") == false {
				break
			} else {
				printMsg("Grapper", "media play list url is chunklist, assume live stream has not already started. Try in 1 minute. "+fmt.Sprintf("%v", maxTries-counter))
				time.Sleep(1 * time.Minute)
			}

			counter++
		}

		printMsg("Grapper", fmt.Sprintf("mediaPlaylistUrl: %v", mediaPlaylistUrl))
		mediaPlaylistUrls <- mediaPlaylistUrl

		// Set timeslot
		timeSlots <- order.timeSlot
	}
}

func MediaLoaderSlidingWindow(mediaPlaylistUrls <-chan cm3u8.M3U8URL, startSignal, stopSignal <-chan bool) {

	intervalSec := uint(10) // TODO: Interval from mediaPlayList
	baseUrl := cm3u8.M3U8URL("")
	onOffSignal := make(chan bool)

	mediaPlaylistUrlsIn := make(chan cm3u8.M3U8URL)
	mediaPlaylistUrlsOut := make(chan cm3u8.M3U8URL)
	mediaPlaylistUrlsSwitched := make(chan cm3u8.M3U8URL)
	mediaPlaylists := make(chan m3u8.MediaPlaylist)
	mediaSegmentsUris := make(chan cm3u8.M3U8URL)

	urlsToDownload := make(chan cm3u8.M3U8URL, 100)
	urlsNotAlreadyDownloaded := make(chan cm3u8.M3U8URL, 100)
	downloadedUrls := make(chan cm3u8.M3U8URL, 100)
	downloadItems := make(chan DownloadItem, 100)

	//          interval         on/off
	// +---------------------------------------------------------------------------
	// |            |               |
	// |   url +----v-----+    +----v---+    +-------------+    +-------------+
	// |   --->| Repeater |--->| Switch |--->| MediaLoader |--->| SegsGrapper |--->
	// |       +----------+    +--------+    +-------------+    +-------------+
	// |
	// |           +--------+    +-----------+
	// |       --->| absUrl |--->| FilterDLD |--->
	// |           +--------+    +-----------+
	// |
	// +----------------------------------------------------------------------------

	go cm3u8.Repeater(intervalSec, mediaPlaylistUrlsIn, mediaPlaylistUrlsOut)
	go cm3u8.Switch(onOffSignal, mediaPlaylistUrlsOut, mediaPlaylistUrlsSwitched)
	// TODO: MAPPER
	go cm3u8.MediaLoader(mediaPlaylistUrlsSwitched, mediaPlaylists)
	// TODO: MAPPER
	go cm3u8.SegmentsGrapper(mediaPlaylists, mediaSegmentsUris)

	// Make Absolute (MAP)
	//
	absUrlFn := func(url cm3u8.M3U8URL) cm3u8.M3U8URL {
		url = makeAbsolute(baseUrl, cm3u8.M3U8URL(url))
		//fmt.Println("MakeAbsolute", "-", "url:", url)
		return url
	}
	go cm3u8.Mapper(mediaSegmentsUris, urlsToDownload, absUrlFn)

	// Filter away already downloaded urls or queued urls (FILTER)
	//
	downloadedUrlsMap := map[cm3u8.M3U8URL]bool{}
	filterFn := func(url cm3u8.M3U8URL) bool {
		_, present := downloadedUrlsMap[url]

		msg := fmt.Sprint("Filter", "-", "url:", url, "present:", present)
		printMsg("Controller", msg)

		return (present == false) // if not in map then download, let it pass through filter
	}
	go cm3u8.Filter(urlsToDownload, urlsNotAlreadyDownloaded, filterFn)

	go func() {
		for {
			url := <-urlsNotAlreadyDownloaded
			//fmt.Println("DownloadItems creator", "-", "url:", url)
			newItem := DownloadItem{
				url:    url,
				folder: "./download/",
			}
			downloadItems <- newItem
		}
	}()

	go Downloader(downloadItems, downloadedUrls)

	go func() {
		for {
			url := <-downloadedUrls
			msg := fmt.Sprint("Downloaded Items", "-", "url:", url)
			printMsg("Controller", msg)
			downloadedUrlsMap[url] = true
		}
	}()

	go func() {
		for {
			select {
			case <-startSignal:
				printMsg("Controller", "START signal received!")
				onOffSignal <- true

			case <-stopSignal:
				printMsg("Controller", "STOP signal received!")
				onOffSignal <- false
			}
		}
	}()

	for {
		url := <-mediaPlaylistUrls
		baseUrl = cm3u8.GetBaseUrl(url)
		mediaPlaylistUrlsIn <- url
	}
}

func main() {

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

	// download what?
	//  - channel ?
	//  - when? start, end
	//  - to folder ?

	dlWhat := func() {

		channel := getString("Which channel" + getChannelList() + "?")
		startTimeUTC := getDateTimeLocal("Which start time?").UTC()
		endTimeUTC := getDateTimeLocal("Which end time?").UTC()
		folder := getString("Which folder?")

		// Check
		// - Channel exists
		// - EndTime > StartTime
		// - Check folder exists else create it

		dlo := DownloadOrder{
			channel: channel,
			timeSlot: TimeSlot{
				start: startTimeUTC,
				end:   endTimeUTC,
			},
			folder: folder,
		}

		if validate(dlo) {
			// Send order into network
			//
			orders <- dlo
		} else {
			printMsg("Application", "Order not valid!")
		}

	}

	cmds := map[string]func(){
		"quit": func() { quit <- true },
		"dl":   func() { dl() },
		"dlw":  func() { dlWhat() },
	}
	go cell.Console(cmds) // stdout, stdin
	go Grapper(orders)

	// Wait until quit
	//
	<-quit
	printMsg("Application", "Quit now!!")
}
