package main

import (
	"bufio"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/grafov/m3u8"
	"github.com/mysheep/cell/cm3u8"
	"github.com/mysheep/cell/ctime"
	"github.com/mysheep/cell/web"
)

const (
	// dateFormat uses in getDatetimeLocal function
	dateFormat     = "2006-01-02 15:04"
	dateFormatLong = "2006-01-02 15:04:05"
)

// see https://stackoverflow.com/questions/25318154/convert-utc-to-local-time-go
// countryTz map of town to location
var countryTz = map[string]string{
	"Vienna": "Europe/Vienna",
	// ...
}

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

// createDownloadOrder create a downloadOrder with the given
// parameter, it also ensure that the folder exists are creation
func createDownloadOrder(channel string, startTimeUTC, endTimeUTC time.Time, folder string) (DownloadOrder, error) {
	dir, err := ensureDir(folder)
	if err != nil {
		return DownloadOrder{}, err
	}
	return DownloadOrder{
		channel: channel,
		timeSlot: TimeSlot{
			start: startTimeUTC,
			end:   endTimeUTC,
		},
		folder: dir,
	}, nil
}

func validate(do DownloadOrder) bool {

	// validateChannel validates, if the entered channel is a supported channel
	validateChannel := func() bool {
		chs := getKeys(channels)
		fmt.Println("chs:", chs)
		if stringInSlice(do.channel, chs) == false {
			printMsg("Validater", "channel '"+do.channel+"' not in list!")
			return false
		}
		return true
	}

	// validateTimeSlot validates, if start time is before end time
	validateTimeSlot := func() bool {
		if do.timeSlot.start.Sub(do.timeSlot.end) > 0 {
			printMsg("Validater", "start time > end time!")
			return false
		}
		return true
	}

	// validateFolder validates, if folder string is a folder
	validateFolder := func() bool {
		_, err := os.Stat(do.folder)
		if err != nil {
			printMsg("Validater", "folder '"+do.folder+"' is not a path!")
			return false
		}
		return true
	}

	return validateChannel() && validateTimeSlot() && validateFolder()
}

// getChannelList return a list of
// the channels e.g. (orf1, orf2)
func getChannelList() string {
	keys := getKeys(channels)
	return "(" + strings.Join(keys, ", ") + ")"
}

// getKeys return the keys of the given
// channels map
func getKeys(channels map[string]cm3u8.M3U8URL) []string {
	keys := make([]string, 0, len(channels))
	for k, _ := range channels {
		keys = append(keys, k)
	}
	return keys
}

// createFullDir creates the full download directory path
// e.g. folder = "bar" -> fullDir = "./download/bar/"
func createFullDir(folder string) string {
	fullpath := filepath.Join(".", downloadFolder, folder, "/")
	return fullpath
}

// existsDir check if a given folder already exists
func existsDir(folder string) bool {
	fileInfo, err := os.Stat(folder)
	if err == nil && fileInfo.IsDir() {
		return true
	}
	return false
}

// ensureDir ensure that a given directory
// exists, if it not already exists it
// create the directory
func ensureDir(folder string) (string, error) {
	err := os.MkdirAll(folder, os.ModePerm)
	if err != nil {
		return "", err
	}
	return folder, nil
}

// stringInSlice return true if the given
// string a in list of given strings
// TODO: General library !!!
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

/*
2019-04-27 14:21:09             Grapper - It seams no livestream available!
2019-04-27 14:21:09             Grapper - It seams no livestream available!
2019-04-27 14:21:09             Grapper - It seams no livestream available!
*/

// printMsg print a debug msg of object to console
func printMsg(object string, msg string) {
	now := time.Now().Format(dateFormatLong)
	fmt.Printf("%s %20s - %s\n", now, object, msg)
}

// getString ask the user to enter a string
func getString(question string) string {
	var result string
	fmt.Print(question)
	fmt.Print(" ")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		result = scanner.Text()
	}
	return result
}

// getDateTimeLocal ask the user to enter a date time
// e.g. > Start time ? 2019-04-01 13:10 and it convert
// to local time (location is "Vienna")
func getDateTimeLocal(question string) time.Time {
	result := time.Now()
	for {
		dateTimeStr := getString(question)

		loc, err := time.LoadLocation(countryTz["Vienna"])
		if err != nil {
			panic(err)
		}

		result, err = time.ParseInLocation(dateFormat, dateTimeStr, loc)
		if err == nil {
			break
		} else {
			fmt.Println("Wrong format. Try again.")
		}
	}
	return result
}

// getDownloadFolder ask the user for the relative folder inside download
// e.g. if user enter "movieA", it returns "download/movie"
func getDownloadSubFolder(question string) string {
	fullpath := createFullDir(".")
	for {
		relFolder := getString(question)
		fullpath = createFullDir(relFolder)
		if existsDir(fullpath) {
			fmt.Println("Folder", fullpath, "already exists!")
		} else {
			break
		}
	}
	return fullpath
}

// getFileName extracts the file name of an url
// e.g. http://foo.com/bar.ts -> bar.ts
func getFilename(urlRaw cm3u8.M3U8URL) string {
	url, err := url.Parse(string(urlRaw))
	if err != nil {
		panic(err)
	}
	return path.Base(url.Path)
}

// makeAbsolute check if the url is a relative
// and then make an absolute url out of it
func makeAbsolute(base, url cm3u8.M3U8URL) cm3u8.M3U8URL {
	if cm3u8.IsRelativeUrl(url) {
		return base + url
	}
	return url
}

// StartStopTimmer get a timeslot with start and stop time
// and send a signal at start time to startSignals channel
// and a signal at end time to stopSignals channel
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

func getMediaPlayListUrlOfVariant(baseUrl cm3u8.M3U8URL, masterPlaylist m3u8.MasterPlaylist, variantIndex uint) (cm3u8.M3U8URL, error) {

	mediaPlaylistUrl := cm3u8.M3U8URL(masterPlaylist.Variants[variantIndex].URI)
	mediaPlaylistUrl = makeAbsolute(baseUrl, mediaPlaylistUrl)

	// https://apasfiis.sf.apa.at/ipad/gp/livestream_Q6A.mp4/chunklist.m3u8?lbs=20190412132743573&origin=http%253a%252f%252fvarorfvod.sf.apa.at%252fsystem_clips%252flivestream_Q6A.mp4%252fchunklist.m3u8&ip=129.27.216.70&ua=Go-http-client%252f1.1

	if strings.Contains(string(mediaPlaylistUrl), "chunklist.m3u8") {
		return cm3u8.M3U8URL(""), errors.New("Media play list is chunklist")
	}

	return mediaPlaylistUrl, nil
}

// Downloader cell download the given url in downloaditem
// and save it to the given folder
func Downloader(items <-chan DownloadItem, downloaded chan<- cm3u8.M3U8URL) {

	// Channels for Downloader
	urls := make(chan string)
	contents := make(chan []byte)

	// Channels for Saver
	filenames := make(chan string)
	bytess := make(chan []byte)
	savedFilenames := make(chan string)

	// Setup SubNet
	go web.Downloader(urls, contents)
	go web.Saver(filenames, bytess, savedFilenames)

	// Now go and work
	for {
		// receive information what to download and save
		item := <-items
		fileName := path.Join(item.folder, getFilename(item.url))

		// Send url to downloader and ...
		urls <- string(item.url)

		// ... wait for downloaded content
		content := <-contents

		// Send fileName and content to saver and ...
		filenames <- fileName
		bytess <- content

		// ... wait for file is saved
		<-savedFilenames
		downloaded <- item.url
	}
}
