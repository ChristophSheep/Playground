package main

import (
	"io/ioutil"
	"net/http"
)

import (
	"fmt"
	"strings"

	"github.com/mysheep/conns"
	"github.com/mysheep/whois"
)

func printConnTable(conns []conns.Connection) {
	for _, conn := range conns {

		// 0.0.0.0:80
		// [IPv4 or v6]:[PORT]
		//
		name_values, err := whois.Get(getIP(conn.Wan))

		name := "empty"
		if err == nil {

			if str, found := name_values["Organization"]; found {
				name = str
			}

			if str, found := name_values["address"]; found {
				name = str
			}
		}

		fmt.Printf("%50s %50s %50s\n", conn.Lan, conn.Wan, name)
	}
}

func printMap(m map[string]string) {
	for k, v := range m {
		fmt.Printf("%20s: %s\n", k, v)
	}
}

func reverse(input string) string {
	n := 0
	rune := make([]rune, len(input))
	for _, r := range input {
		rune[n] = r
		n++
	}
	rune = rune[0:n]
	// Reverse
	for i := 0; i < n/2; i++ {
		rune[i], rune[n-1-i] = rune[n-1-i], rune[i]
	}
	// Convert back to UTF-8.
	return string(rune)
}

func getIP(ipPort string) string {
	xs := strings.SplitAfterN(reverse(ipPort), ":", 2)
	ip := xs[1]
	ip = reverse(ip)
	return strings.TrimRight(strings.TrimLeft(ip, "["), "]")
}

/*
func main() {

	cs, err := conns.Get()

	if err != nil {
		fmt.Print(err)
	}

	printConnTable(cs)

}
*/
var image []byte

// preparing image
func init() {
	var err error
	image, err = ioutil.ReadFile("./image.png")
	if err != nil {
		panic(err)
	}
}

// Send HTML and push image
func handlerHtml(w http.ResponseWriter, r *http.Request) {
	pusher, ok := w.(http.Pusher)
	if ok {
		fmt.Println("Push /image")
		pusher.Push("/image", nil)
	}
	w.Header().Add("Content-Type", "text/html")
	fmt.Fprintf(w, `<html><body><img src="/image"></body></html>`)
}

// Send image as usual HTTP request
func handlerImage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/png")
	w.Write(image)
}
func main() {
	http.HandleFunc("/", handlerHtml)
	http.HandleFunc("/image", handlerImage)
	fmt.Println("start http listening :18443")
	//err := http.ListenAndServeTLS(":18443", "server.crt", "server.key", nil)
	err := http.ListenAndServe(":18443", nil)
	fmt.Println(err)
}
