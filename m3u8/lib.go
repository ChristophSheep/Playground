package main

import (
	"fmt"
	"net/url"
	"path"

	"github.com/mysheep/cell/cm3u8"
)

func printMsg(object string, msg string) {
	fmt.Printf("%25s - %s\n", object, msg)
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
