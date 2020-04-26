/*
"use strict";

const WebSocketServer = require("ws").Server;
const wss = new WebSocketServer({ port: 2222 });

wss.on("connection", (ws) => {
   console.info("websocket connection open");

   if (ws.readyState === ws.OPEN) {
       ws.send(JSON.stringify({
           msg1: 'yo, im msg 1'
       }))

       setTimeout(() => {
            ws.send(JSON.stringify({
                msg2: 'yo, im a delayed msg 2'
            }))
       }, 1000)
   }
});
*/

package main

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

var address = "localhost:2222"
var origin = "http://localhost/"
var fooUrl = "ws://localhost:2222/foo"
var echoUrl = "ws://localhost:2222/echo"

func echoHandler(ws *websocket.Conn) {

	for {
		_, err := ws.Write([]byte("hello"))
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func fooHandler(ws *websocket.Conn) {

	for {
		_, err := ws.Write([]byte("foo"))
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(25 * time.Millisecond)
	}
}

func createServer(address string) {

	// Register Handlers
	http.Handle("/echo", websocket.Handler(echoHandler))
	http.Handle("/foo", websocket.Handler(fooHandler))

	// Listen to clients connecting
	err := http.ListenAndServe(address, nil)

	// Show Error
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}

}
func main2() {
	createServer(address)
}
