package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/sac007/gowebsocket"
)

func main() {

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	host := "127.0.0.1"
	port := "7777"
	schema := "ws"

	// init
	// schema – can be ws or wss
	// host, port – ws server
	//socket := gowebsocket.New({schema}://{host}:{port})

	socket := gowebsocket.New(fmt.Sprintf("%s://%s:%s/A", schema, host, port))

	socket.OnConnected = func(socket gowebsocket.Socket) {
		log.Println("Connected to server")
	}

	socket.OnConnectError = func(err error, socket gowebsocket.Socket) {
		log.Println("Recieved connect error ", err)
	}

	// send message
	// socket.SendText({message})

	//or

	//socket.SendBinary({message})

	// receive message
	socket.OnTextMessage = func(message string, socket gowebsocket.Socket) {
		// hande received message
		fmt.Printf("Received message %s from server", message)
	}

	socket.Connect() // to server

	for {
		select {
		case <-interrupt:
			log.Println("interrupt")
			socket.Close()
			return
		}
	}
}
