package main

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/grafov/m3u8"
	"github.com/mysheep/cell"
	"github.com/mysheep/cell/cm3u8"
)

const (
	DEBUG          = false
	downloadFolder = "./download/"
)

var (
	channels = map[string]cm3u8.M3U8URL{
		"orf1": cm3u8.M3U8URL("http://orf1.orfstg.cdn.ors.at/out/u/orf1/q6a/manifest.m3u8"), // hoch
		"orf2": cm3u8.M3U8URL("http://orf2.orfstg.cdn.ors.at/out/u/orf2/q4a/manifest.m3u8"), // mittel
	}
)

func Grapper(orders <-chan DownloadOrder) {

	timeSlots := make(chan TimeSlot)
	startSignal := make(chan bool)
	stopSignal := make(chan bool)
	folders := make(chan string)

	masterUrls := make(chan cm3u8.M3U8URL)

	// Setup Network
	//
	go StartStopTimer(timeSlots, startSignal, stopSignal)
	go MediaLoaderSlidingWindow(masterUrls, startSignal, stopSignal, folders)

	// Get the order
	//
	for {

		order := <-orders
		printMsg("Grapper", fmt.Sprint("order: ", order))
		masterUrl := channels[order.channel]
		printMsg("Grapper", fmt.Sprintf("masterUrl: %s", masterUrl))

		// Set master url to download ..
		masterUrls <- masterUrl
		folders <- order.folder
		// Set timeslot
		timeSlots <- order.timeSlot
	}
}

func MediaLoaderSlidingWindow(masterUrlsIn <-chan cm3u8.M3U8URL, startSignal, stopSignal <-chan bool, folders <-chan string) {

	intervalSec := uint(10) // TODO: Interval from mediaPlayList
	baseUrl := cm3u8.M3U8URL("")
	onOffSignal := make(chan bool)

	// Repeater
	masterPlaylistUrlsIn := make(chan cm3u8.M3U8URL)
	masterPlaylistUrlsRepeated := make(chan cm3u8.M3U8URL)

	// MasterPlayListLoader
	masterPlaylists := make(chan m3u8.MasterPlaylist)

	mediaPlaylistUrls := make(chan cm3u8.M3U8URL)

	mediaPlaylistUrlsSwitched := make(chan cm3u8.M3U8URL)
	mediaPlaylists := make(chan m3u8.MediaPlaylist)
	mediaSegmentsUris := make(chan cm3u8.M3U8URL)

	urlsToDownload := make(chan cm3u8.M3U8URL, 100)
	urlsNotAlreadyDownloaded := make(chan cm3u8.M3U8URL, 100)
	downloadedUrls := make(chan cm3u8.M3U8URL, 100)
	downloadItems := make(chan DownloadItem, 100)

	folder := downloadFolder // Default ./download/
	go func() {
		folder = <-folders
	}()

	//  folder     interval                       on/off
	//     |         |                              |
	//  +--+---------+------------------------------+-------------------------------+
	//  |            |                              |                               |
	//  |   url +----v-----+   +-----------+   +----v---+   +----------+ 		    |
	// -+------>| Repeater |-->| MasterLdr |-->| Switch |-->| MediaLdr |--> ...     |
	//  |       +----------+   +-----------+   +--------+   +----------+            |
	//  |																		    |
	//  |          +------------+   +--------+   +-----------+   +------------+     |
	//  |   ... -->|SegsGrapper |-->| absUrl |-->| FilterDLD |-->| Downloader |-->  |
	//  |          +------------+   +--------+   +-----------+   +------------+     |
	//  | 																		    |
	//  +---------------------------------------------------------------------------+
	//

	go cm3u8.Repeater(intervalSec, masterPlaylistUrlsIn, masterPlaylistUrlsRepeated)
	go cm3u8.MasterLoader(masterPlaylistUrlsRepeated, masterPlaylists)

	go func() {

		getMediaPlayListUrl := func(masterPlaylist m3u8.MasterPlaylist) (cm3u8.M3U8URL, error) {

			// TODO: Take first variant or ??
			//
			mediaPlaylistUrl := cm3u8.M3U8URL(masterPlaylist.Variants[0].URI)
			mediaPlaylistUrl = makeAbsolute(baseUrl, mediaPlaylistUrl)

			// https://apasfiis.sf.apa.at/ipad/gp/livestream_Q6A.mp4/chunklist.m3u8?lbs=20190412132743573&origin=http%253a%252f%252fvarorfvod.sf.apa.at%252fsystem_clips%252flivestream_Q6A.mp4%252fchunklist.m3u8&ip=129.27.216.70&ua=Go-http-client%252f1.1

			if strings.Contains(string(mediaPlaylistUrl), "chunklist.m3u8") {
				return cm3u8.M3U8URL(""), errors.New("Media play list is chunklist")
			}

			return mediaPlaylistUrl, nil
		}

		for {
			masterPlaylist := <-masterPlaylists // from Repeater

			mediaPlaylistUrl, err := getMediaPlayListUrl(masterPlaylist)
			if err == nil {
				//printMsg("Grapper", fmt.Sprintf("mediaPlaylistUrl: %v", mediaPlaylistUrl))
				mediaPlaylistUrls <- mediaPlaylistUrl
			} else {
				printMsg("Grapper", "It seams no livestream aavailable!")
			}
		}
	}()

	go cm3u8.Switch(onOffSignal, mediaPlaylistUrls, mediaPlaylistUrlsSwitched)
	go cm3u8.MediaLoader(mediaPlaylistUrlsSwitched, mediaPlaylists)
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

		//msg := fmt.Sprint("Filter", "-", "url:", url, "present:", present)
		//printMsg("Controller", msg)

		return (present == false) // if not in map then download, let it pass through filter
	}
	go cm3u8.Filter(urlsToDownload, urlsNotAlreadyDownloaded, filterFn)

	go func() {
		for {
			url := <-urlsNotAlreadyDownloaded
			//fmt.Println("DownloadItems creator", "-", "url:", url)
			newItem := DownloadItem{
				url:    url,
				folder: folder,
			}
			downloadItems <- newItem
		}
	}()

	go Downloader(downloadItems, downloadedUrls)

	go func() {
		for {
			url := <-downloadedUrls
			msg := fmt.Sprint("Downloaded item:", getFilename(url))
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
		masterUrl := <-masterUrlsIn
		baseUrl = cm3u8.GetBaseUrl(masterUrl)
		masterPlaylistUrlsIn <- masterUrl // send to Repeater
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

	// dlw is short cut for download what?
	// it ask the user for channel, start and end time
	// and folder to save the live stream files
	dlw := func() {

		channel := getString("Which channel " + getChannelList() + " ?")
		startTimeUTC := getDateTimeLocal("Which start time (eg. 2019-04-04 12:00) ?").UTC()
		endTimeUTC := getDateTimeLocal("Which   end time (eg. 2019-04-04 12:00) ?").UTC()
		folder := getString("Which folder ?")

		// Check
		// - Channel exists
		// - EndTime > StartTime
		// - Check folder exists else create it

		dlo, err := createDownloadOrder(
			channel,
			startTimeUTC,
			endTimeUTC,
			folder,
		)

		if validate(dlo) && err == nil {
			// Send order into network
			orders <- dlo
		} else {
			printMsg("Application", "Order not valid! No order send into network!")
		}

	}

	// Map of current builtin commands
	//
	cmds := map[string]func(){
		"quit": func() { quit <- true },
		"dl":   func() { dl() },
		"dlw":  func() { dlw() },
	}
	go cell.Console(cmds) // stdout, stdin
	go Grapper(orders)

	// Wait until quit
	//
	<-quit
	printMsg("Application", "Quit now!!")
}
