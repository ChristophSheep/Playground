package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"path"
	"time"

	"github.com/mysheep/cell/cm3u8"
	"github.com/mysheep/cell/ctime"
	"github.com/mysheep/cell/web"
)

const (
	dateFormat = "2006-01-02 15:04"
)

// see https://stackoverflow.com/questions/25318154/convert-utc-to-local-time-go

var countryTz = map[string]string{
	"Vienna": "Europe/Vienna",
	// ...
}

func printMsg(object string, msg string) {
	fmt.Printf("%25s - %s\n", object, msg)
}

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

func getDateTime(question string) time.Time {
	dateTimeStr := getString(question)

	loc, err := time.LoadLocation(countryTz["Vienna"])
	if err != nil {
		panic(err)
	}

	result, err := time.ParseInLocation(dateFormat, dateTimeStr, loc)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	return result
}

func getFilename(urlRaw cm3u8.M3U8URL) string {

	url, err := url.Parse(string(urlRaw))
	if err != nil {
		panic(err)
	}

	return path.Base(url.Path)
}

func makeAbsolute(base, url cm3u8.M3U8URL) cm3u8.M3U8URL {
	if cm3u8.IsRelativeUrl(url) {
		return base + url
	}
	return url
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
