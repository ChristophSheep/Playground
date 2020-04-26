package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

var origin = "http://localhost/"
var url = "ws://localhost:8080/echo"

func inputHandler(ws *websocket.Conn) {

	// Message lesen
	msg := make([]byte, 512)
	// Read 512 bytes
	n, err := ws.Read(msg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Receive: %s\n", msg[:n])

}

func connectTo() {
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}

	message := []byte("hello, world!")
	_, err = ws.Write(message)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Send: %s\n", message)
}

func main() {

	getPort := func() int {

		// see https://gobyexample.com/command-line-subcommands
		// node -port=8001
		//

		portPtr := flag.Int("port", 8080, "an int")
		flag.Parse()

		fmt.Printf("Try to start server at port: %v\n", *portPtr)
		return *portPtr
	}

	getAddr := func(port int) string {
		addr := fmt.Sprintf(":%v", port)
		return addr
	}

	//
	// [Client]------>[Server  Client]------>[Server]
	//

	port := getPort()
	addr := getAddr(port)

	http.Handle("/input", websocket.Handler(inputHandler))

	err := http.ListenAndServe(addr, nil) // Waits hier
	if err != nil {
		panic("Error ListenAndServe: " + err.Error())
	}

}
