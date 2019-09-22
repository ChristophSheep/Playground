package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

const (
	port       = "8080"
	rootFolder = "/Users/christophreif/Movies/F1/2019/"
)

//
// see https://wiki.selfhtml.org/wiki/CSS/Tutorials/Bilder_mit_Bildunterschriften
//
const tpl = `
{{define "Item"}}
<div>{{.}}</div>
<figure>
	<a href="./F1/2019/{{.Folder}}/master.m3u8">
		<img src="./F1/images/{{.Name}}.png">
		<figcaption>{{.Date}} {{.Name}}</figcaption>
	</a>
</figure>
{{end}}

<!DOCTYPE html>
<html>
    <head>
        <meta charset="UTF-8">
		<title>{{.Title}}</title>
		<style>
			figure {
				position: relative;
				margin: 20px;
				padding: 10px;
				border: 1px solid gainsboro;
				background: white;
			}
			
			figcaption {
				padding: 10px;
				text-align: center;
			}
		</style>
    </head>
	<body>
		<h1>{{.Title}}</h1>
		{{range .Items}}
		{{template "Item" .}}
		{{end}}
    </body>
</html>`

// http://localhost:8000/F1/2019/2019-04-13_ChinaQP/master.m3u8
// http://localhost:8000/F1/Images/Baku.png

func getFolders(parentFolder string) []string {
	file, err := os.Open(parentFolder)
	check(err)

	dirs, err := file.Readdirnames(-1)
	check(err)

	filtered := []string{}

	for _, dir := range dirs {
		fi, err := os.Stat(path.Join(rootFolder, dir))
		//check(err)
		if err == nil {
			if fi.IsDir() {
				filtered = append(filtered, dir)
			}
		}
	}

	return filtered
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type Page struct {
	Title string
	Items []Item
}

type Item struct {
	Folder string // Attribute must start with BIG letter otherwise its private
	Date   string
	Name   string
	Kind   string
}

// getDate get the date from folder name
// e.g. 2019-05-12_Spanien-Barcelona_FT -> 2019-05-12
func getDate(folder string) string {
	return strings.Split(folder, "_")[0]
}

// getName get the name from the the folder name
// e.g. 2019-05-12_Spanien-Barcelona_FT -> Spanien-Barcelona
func getName(folder string) string {
	return strings.Split(folder, "_")[1]
}

// getName get the name from the the folder name
// e.g. 2019-05-12_Spanien-Barcelona_FT -> Spanien-Barcelona
func getType(folder string) string {
	return strings.Split(folder, "_")[2]
}

func createItem(folder string) Item {
	return Item{
		Folder: folder,
		Date:   getDate(folder),
		Name:   getName(folder),
		Kind:   getType(folder),
	}
}

func mapItem(folders []string) []Item {
	items := make([]Item, len(folders))
	for i, folder := range folders {
		item := createItem(folder)
		items[i] = item
	}
	return items
}

func writeFolders(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("request: %v", r)

	// Create template for HTML page
	//
	t, err := template.New("webpage").Parse(tpl)
	check(err)

	// Get folders
	//
	folders := getFolders(rootFolder)

	// Create data
	//
	data := Page{
		Title: "F1 Movies",
		Items: mapItem(folders),
	}

	// Render View by Template
	//
	err = t.Execute(w, data)
}

func main() {
	http.HandleFunc("/", writeFolders)
	fmt.Println("start server at port", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}
