package main

import (
	"bufio"
	"errors"
	"fmt"
	"html/template"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const masterPlayListFileName = "master.m3u8"
const mediaPlayListFileName = "media.m3u8"

const masterT = `#EXTM3U
#EXT-X-VERSION:4
#EXT-X-STREAM-INF:BANDWIDTH=2194075,RESOLUTION=960x540,CODECS="avc1.4D401F,mp4a.40.2"
{{.MediaUrl}}
`

//
// #EXT-PLAYLIST-TYPE:VOD   // Not LiveStream
// #EXT-PLAYLIST-TYPE:EVENT // LiveStream

// see https://golang.org/pkg/text/template/
const mediaT = `#EXTM3U
#EXT-X-VERSION:4
#EXT-PLAYLIST-TYPE:VOD
#EXT-X-TARGETDURATION:10
#EXT-X-MEDIA-SEQUENCE:{{.ExtXMediaSequence}}
{{range .Segments}}
#EXTINF:10.000, 
{{.URI}}
{{end}}
#EXT-X-ENDLIST
`

type Master struct {
	MediaUrl string
}

type Media struct {
	ExtXMediaSequence string
	Segments          []Segment
}

type Segment struct {
	URI string
}

func createMasterFile(folder string) {
	//define an instance
	data := Master{
		MediaUrl: mediaPlayListFileName,
	}

	template := getMasterTemplate()
	fullPath := path.Join(folder, masterPlayListFileName)
	writeTemplate(fullPath, template, data)
}

func getMasterTemplate() template.Template {
	//create a new template with some name
	masterTempl := template.New("MASTERPL")
	//parse some content and generate a template
	tmpl, err := masterTempl.Parse(masterT)
	if err != nil {
		log.Fatal("Parse: ", err)
		panic(err)
	}
	return *tmpl
}

func getMediaTemplate() template.Template {
	//create a new template with some name
	masterTempl := template.New("MEDIAPL")
	//parse some content and generate a template
	tmpl, err := masterTempl.Parse(mediaT)
	if err != nil {
		log.Fatal("Parse: ", err)
		panic(err)
	}
	return *tmpl
}

func getStartSegNumber(segments []Segment) (string, error) {

	if len(segments) == 0 {
		return "", errors.New("Segments length is 0")
	}

	first := segments[0]
	fileName := path.Base(first.URI) // e.g. manifest_4_123.ts

	ext := path.Ext(fileName)                    // e.g. .ts
	fileName = strings.TrimSuffix(fileName, ext) // e.g. mainfest_4_123

	isDigit := func(c byte) bool {
		return c >= '0' && c <= '9'
	}

	num := ""
	for i := len(fileName) - 1; i >= 0; i-- {
		c := fileName[i]
		if isDigit(c) {
			num = string(c) + num
		} else {
			break
		}
	}

	return num, nil
}
func createMediaFile(folder string, tsFiles []string) {

	createSegments := func(tsFiles []string) []Segment {
		var segments []Segment
		for _, tsFile := range tsFiles {
			s := Segment{
				URI: tsFile,
			}
			segments = append(segments, s)
		}
		return segments
	}

	segments := createSegments(tsFiles)
	fmt.Printf("Count segments: %v\n", len(segments))

	startNum, err := getStartSegNumber(segments)
	if err != nil {
		fmt.Println("No sgements")
	} else {
		fmt.Println("Start seq number found:", startNum)
	}

	data := Media{
		ExtXMediaSequence: startNum,
		Segments:          segments,
	}

	template := getMediaTemplate()
	fullPath := path.Join(folder, mediaPlayListFileName)
	writeTemplate(fullPath, template, data)

}

func writeTemplate(fullPath string, tmpl template.Template, data interface{}) {
	// create new file
	fo, err := os.Create(fullPath)
	if err != nil {
		panic(err)
	}
	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

	// Create a writer to write to file
	w := bufio.NewWriter(fo)

	// merge template 'tmpl' with content of 's'
	// and write to file
	err1 := tmpl.Execute(w, data)
	if err1 != nil {
		log.Fatal("Execute: ", err1)
		return
	}

	// Flush writer
	if err = w.Flush(); err != nil {
		panic(err)
	}
}

func getTsFiles(folder string) []string {
	var filepaths []string

	filter := func(filepath string) bool {
		return path.Ext(filepath) == ".ts"
	}

	root := folder
	filepath.Walk(root, func(filepath string, info os.FileInfo, err error) error {
		if filter(filepath) {
			if path.IsAbs(filepath) {
				filepath = path.Base(filepath) // Base returns last element of path
			}
			filepaths = append(filepaths, filepath)
		}
		return nil
	})
	return filepaths
}

/*
func writeJnSh(script string) {
	fo, _ := os.Create("jn.sh")
	defer fo.Close()
	w := bufio.NewWriter(fo)
	w.WriteString(script)
	w.Flush()
}
*/
/*
func main() {
	const folder = "." // Current folder
	fileNames := getTsFiles(folder)
	script := "cat " + strings.Join(fileNames, " ") + " > movie.ts"
	fmt.Println(script)
	writeJnSh(script)
}
*/

func main() {
	const folder = "/Users/christophreif/Movies/Formel1/2019-04-14_ChinaGP"
	//const folder = "."

	files := getTsFiles(folder)
	createMasterFile(folder)
	createMediaFile(folder, files)
}
