package main

import (
	"fmt"
	"time"

	"golang.org/x/net/websocket"
)

var origin = "http://localhost/"
var url = "ws://localhost:2222/echo"

func main2() {
	ws, err := websocket.Dial(url, "", origin)

	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		var m string
		websocket.Message.Receive(ws, &m)
		//websocket.Handler.ServeHTTP()

		fmt.Println(m)
		time.Sleep(500 * time.Millisecond)
	}
}
