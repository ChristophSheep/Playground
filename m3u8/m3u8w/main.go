package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
)

const (
	port       = "8080"
	rootFolder = "/Users/christophreif/Movies/F1/2019/"
)

//
// see https://wiki.selfhtml.org/wiki/CSS/Tutorials/Bilder_mit_Bildunterschriften
//
const tpl = `
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
				width: 470px;
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
		<figure>
			<a href="http://localhost:8000/F1/2019/{{.}}/master.m3u8">
		    	<img src="https://upload.wikimedia.org/wikipedia/commons/thumb/2/29/Korea_international_circuit_v3.svg/260px-Korea_international_circuit_v3.svg.png">
				<figcaption>{{.}}</figcaption>
			</a>
		</figure>
		{{else}}<div><strong>no rows</strong></div>
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

func writeFolders(w http.ResponseWriter, r *http.Request) {
	// Create template for HTML
	//
	t, err := template.New("webpage").Parse(tpl)
	check(err)

	// Get folders
	//
	folders := getFolders(rootFolder)
	//check(err)

	// Create data
	//
	data := struct {
		Title string
		Items []string
	}{
		Title: "F1 Movies",
		Items: folders,
	}

	// Render HTML
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
