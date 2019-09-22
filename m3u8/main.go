package main

import (
	"fmt"
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
	 	"orf1m": cm3u8.M3U8URL("https://orf1.mdn.ors.at/out/u/orf1/q4a/manifest.m3u8"),     // mittel NEU
		"orf1h": cm3u8.M3U8URL("https://orf1.mdn.ors.at/out/u/orf1/q6a/manifest.m3u8"),     // hoch NEUs
		"orf2m": cm3u8.M3U8URL("http://orf2.mdn.ors.at/out/u/orf2/q4a/manifest.m3u8"), 		// mittel
		"orf2h": cm3u8.M3U8URL("http://orf2.mdn.ors.at/out/u/orf2/q6a/manifest.m3u8"), 		// hoch
	}
)

// Grapper graps given download order with master m3u8 url, start, stop time and folder
func Grapper(orders <-chan DownloadOrder) {

	timeSlots := make(chan TimeSlot)
	startSignal := make(chan bool)
	stopSignal := make(chan bool)
	folders := make(chan string)

	masterUrls := make(chan cm3u8.M3U8URL)
	downloadedUrlsOut := make(chan cm3u8.M3U8URL)

	// Setup Network
	//
	go StartStopTimer(timeSlots, startSignal, stopSignal)
	go MediaLoaderSlidingWindow(folders, masterUrls, startSignal, stopSignal, downloadedUrlsOut)

	// Display
	go func() {
		for {
			downloadedUrl := <-downloadedUrlsOut
			printMsg("Grapper", fmt.Sprintf("downloaded url: %s", downloadedUrl))
		}
	}()

	// Get the order
	//
	for {

		order := <-orders
		printMsg("Grapper", fmt.Sprint("order: ", order))
		masterUrl := channels[order.channel]
		// Order is IMPORTANT!!!
		// First send folder then masterUrl
		folders <- order.folder
		masterUrls <- masterUrl
		timeSlots <- order.timeSlot
	}
}

// MediaLoaderSlidingWindow loads segments (*.ts files) of sliding window live stream
func MediaLoaderSlidingWindow(foldersIn <-chan string, masterUrlsIn <-chan cm3u8.M3U8URL,
	startSignalIn, stopSignalIn <-chan bool, downloadedUrlsOut chan<- cm3u8.M3U8URL) {

	// Const
	//
	const intervalSec = uint(10) // TODO: Interval from mediaPlayList

	// Variables
	//
	downloadedUrlsMap := map[cm3u8.M3U8URL]bool{}
	folder := downloadFolder // default

	// SubNet Channels
	//
	onOffSignal := make(chan bool)
	masterPlaylistUrlsIn := make(chan cm3u8.M3U8URL)
	mediaPlaylistUrls := make(chan cm3u8.M3U8URL)
	mediaPlaylistUrlsSwitched := make(chan cm3u8.M3U8URL)
	urlsToDownload := make(chan cm3u8.M3U8URL, 100)
	downloadedUrls := make(chan cm3u8.M3U8URL, 100)

	// SubNet Components
	//

	//  folders           url
	//     |               |
	//  +--+---------------+--------------------
	//  |  |               |
	//  |  folder          |
	//  | 			       |
	//  | intv  +----------v----------+
	// -+------>| RepeatedMediaUrlGrp |
	//  |       +----------+----------+
	//  |                  |
	//  | on/off	  +----v---+
	// -+------------>| Switch |
	// 	|	     	  +----+---+
	//  |                  |
	// 	|		  +--------v--------+
	// 	|		  | SegmentsGrapper |
	// 	|		  +--------+--------+
	// 	|			       |
	//  |   	 +---------v---------+
	//  |        | OnlyNewDownloader |
	//  |   	 +---------+---------+
	// 	|			       |
	//  |   	 +---------v---------+
	//  |        | MarkDownloadedUrl |
	//  |   	 +---------+---------+
	//  |                  |
	//  +------------------+------------
	//                     |

	go RepeatedMediaUrlGrapper(intervalSec, masterPlaylistUrlsIn, mediaPlaylistUrls)
	go cm3u8.Switch(onOffSignal, mediaPlaylistUrls, mediaPlaylistUrlsSwitched)
	go SegmentsGrapper(mediaPlaylistUrlsSwitched, urlsToDownload)

	// Filter away already downloaded urls or queued urls (FILTER)
	//
	filterFn := func(url cm3u8.M3U8URL) bool {
		_, present := downloadedUrlsMap[url]
		//fmt.Println("DOWNLOAD FILTER url:", url, "present:", present)
		return (present == false) // if not in map then download, let it pass through filter
	}
	go OnlyNewDownloader(filterFn, &folder, urlsToDownload, downloadedUrls)
	go StartStopConverter(startSignalIn, stopSignalIn, onOffSignal)

	// Mark downloaded url in map (= queue)
	//
	fn := func(url cm3u8.M3U8URL) cm3u8.M3U8URL {
		downloadedUrlsMap[url] = true
		return url
	}
	go cm3u8.Mapper(downloadedUrls, downloadedUrlsOut, fn)

	for {
		// First set folder To DOWNLOAD IN ..
		folder = <-foldersIn
		// .. then receive the master url
		masterUrl := <-masterUrlsIn
		// .. then send to repeater and let it run
		masterPlaylistUrlsIn <- masterUrl
	}
}

// StartStopConverter converts an incoming startSignal to onOffSignal = true
// and an incoming stopSignal to onOffSignal = false
func StartStopConverter(startSignal, stopSignal <-chan bool, onOffSignal chan<- bool) {
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
}

// SegmentsGrapper graps segments of media playlist from given url
func SegmentsGrapper(urlsIn <-chan cm3u8.M3U8URL, urlsOut chan<- cm3u8.M3U8URL) {
	// Variables
	//
	baseUrl := cm3u8.M3U8URL("")

	// SubNet Channels
	//
	mediaPlaylistUrls := make(chan cm3u8.M3U8URL)
	mediaPlaylists := make(chan m3u8.MediaPlaylist)
	mediaSegmentsUris := make(chan cm3u8.M3U8URL, 100)
	absoluteSegmentsUris := make(chan cm3u8.M3U8URL, 100)

	// SubNet Components
	//
	go cm3u8.MediaLoader(mediaPlaylistUrls, mediaPlaylists)
	go cm3u8.SegmentsGrapper(mediaPlaylists, mediaSegmentsUris)

	// Make Absolute (MAP)
	//
	absUrlFn := func(url cm3u8.M3U8URL) cm3u8.M3U8URL {
		if string(baseUrl) == "" {
			panic("baseUrl is empty")
		}
		url = makeAbsolute(baseUrl, cm3u8.M3U8URL(url))
		return url
	}
	go cm3u8.Mapper(mediaSegmentsUris, absoluteSegmentsUris, absUrlFn)

	go func() {
		for {
			res := <-absoluteSegmentsUris
			urlsOut <- res
		}
	}()

	for {
		in := <-urlsIn
		baseUrl = cm3u8.GetBaseUrl(in)
		mediaPlaylistUrls <- in
	}
}

// RepeatedMediaUrlGrapper graps from master url the media repeated with given interval
func RepeatedMediaUrlGrapper(intervalSec uint, masterUrlsIn <-chan cm3u8.M3U8URL, mediaPlaylistUrlsOut chan<- cm3u8.M3U8URL) {

	// Variables
	//
	baseUrl := cm3u8.M3U8URL("")

	// Constants
	//
	const variantIndex = 0

	// Subnet Channels
	//
	masterPlaylistUrlsIn := make(chan cm3u8.M3U8URL)
	masterPlaylistUrlsRepeated := make(chan cm3u8.M3U8URL)
	masterPlaylists := make(chan m3u8.MasterPlaylist)

	// Subnet Components
	//
	go cm3u8.Repeater(intervalSec, masterPlaylistUrlsIn, masterPlaylistUrlsRepeated)
	go cm3u8.MasterLoader(masterPlaylistUrlsRepeated, masterPlaylists)

	go func() {

		for {
			masterPlaylist := <-masterPlaylists // from Repeater
			mediaPlaylistUrl, err := getMediaPlayListUrlOfVariant(baseUrl, masterPlaylist, variantIndex)
			if err == nil {
				//printMsg("Grapper", fmt.Sprintf("mediaPlaylistUrl: %v", mediaPlaylistUrl))
				mediaPlaylistUrlsOut <- mediaPlaylistUrl
			} else {
				printMsg("Grapper", "It seams no livestream available!")
			}
		}
	}()

	for {
		masterUrl := <-masterUrlsIn
		baseUrl = cm3u8.GetBaseUrl(masterUrl)
		masterPlaylistUrlsIn <- masterUrl // send to Repeater
	}
}

// OnlyNewDownloader only download new url,
// because it check with filterFn if url is already queued or downloaded
// then if skips this url and does not this url to the downloader
func OnlyNewDownloader(filterFn func(url cm3u8.M3U8URL) bool, folderRef *string, // TODO: 'const string'-channel
	urlsIn <-chan cm3u8.M3U8URL, urlsOut chan<- cm3u8.M3U8URL) {

	urlsToDownload := make(chan cm3u8.M3U8URL, 100)
	urlsNotAlreadyDownloaded := make(chan cm3u8.M3U8URL, 100)
	downloadItems := make(chan DownloadItem, 100)
	downloadedUrls := make(chan cm3u8.M3U8URL, 100)

	go cm3u8.Filter(urlsToDownload, urlsNotAlreadyDownloaded, filterFn)

	// DownloadItemCreator
	go func() {
		for {
			url := <-urlsNotAlreadyDownloaded
			newItem := DownloadItem{
				url:    url,
				folder: *folderRef, // IF folder would be a const channel, I do not need type DownloadItem
			}
			downloadItems <- newItem
		}
	}()

	go Downloader(downloadItems, downloadedUrls)

	// Send to OUT
	go func() {
		for {
			durl := <-downloadedUrls
			urlsOut <- durl
		}
	}()

	// Receive from IN
	for {
		url := <-urlsIn
		urlsToDownload <- url
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

	// download what?
	//  - channel ?
	//  - when? start, end
	//  - to folder ?

	// dlw is short cut for download what?
	// it ask the user for channel, start and end time
	// and folder to save the live stream files
	dlw := func() {

		channel := getString("Which channel " + getChannelList() + " ?")
		now := time.Now().Format(dateFormat)
		startTimeUTC := getDateTimeLocal("Which start time (eg.: " + now + ") ?").UTC()
		endTimeUTC := getDateTimeLocal("Which   end time (eg.: " + now + ") ?").UTC()
		folder := getDownloadSubFolder("Which folder ?")

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

	dln := func() {

		channel := getString("Which channel " + getChannelList() + " ?")
		startTimeUTC := time.Now().UTC()
		endTimeUTC := startTimeUTC.Add(24 * time.Hour)
		folder := getDownloadSubFolder("Which folder ?")

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
		"dlw":  func() { dlw() },
		"dln":  func() { dln() },
	}
	go cell.Console(cmds) // stdout, stdin
	go Grapper(orders)

	// Wait until quit
	//
	<-quit
	printMsg("Application", "Quit now!!")
}
