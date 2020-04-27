package wscell

import (
	"fmt"
	"net/http"
	"time"

	"github.com/mysheep/attribute"
	"golang.org/x/net/websocket"
)

// Connection per AttributeName to webSocket.Conn
var connectionMap map[string]*websocket.Conn = map[string]*websocket.Conn{}

// Spec to create a cell
type Spec struct {
	Name        string
	IP          string
	Port        string
	Attributes  []attribute.Attribute
	Connections []Connection
}

// createInputHandler
func createInputHandler(inputName string) func(*websocket.Conn) {
	inputHandler := func(ws *websocket.Conn) {
		for {

			var message string
			err := websocket.Message.Receive(ws, &message)

			now := time.Now().Format(time.StampMilli)
			if err != nil {
				fmt.Printf("%s - Received error '%v'\n", now, err)
			} else {
				fmt.Printf("%s - Received message '%s' at input '%s'\n", now, message, inputName)
			}

			// Wait
			time.Sleep(200 * time.Millisecond)
		}
	}
	return inputHandler
}

func createInputs(spec Spec) {

	// Create an input connection point
	// for each attributes
	// ws://{ip}:{port}/{Name}
	// e.g.
	//   ws://localhost:1234/A
	for _, attr := range spec.Attributes {
		fmt.Printf("Create input connection point at 'ws://%s:%s/%s'\n", spec.IP, spec.Port, attr.Name())
		http.Handle("/"+attr.Name(), websocket.Handler(createInputHandler(attr.Name())))
	}
}

func createAddress(spec Spec) string {
	address := spec.IP + ":" + spec.Port
	return address
}

// CreateAndListen TODO
//
//           *CELL1*					   *CELL2*
//       +------------+					+------------+
//	  -->o /inA  outA o---->  ....  --->o /inA  outA o---->
//	     |            | 				|            |
//	  -->o /inB  outB o---->  ....  --->o /inB  outB o---->
//       +------------+					+------------+
//         [IP]:[PORT]		webSocket	  [IP]:[PORT]
//
//    server        client			 server        client
//
func CreateAndListen(spec Spec) {

	done := make(chan bool)

	var err error

	// Create INPUT connection points
	// e.g. can be connect via "ws://localhost:1234/inputA"
	createInputs(spec)

	// Create OUTPUT connections
	createOutputConnections(spec)

	// Listen to clients connecting
	// ws://{ip}:{port}/{Name}
	// e.g.
	//   ws://localhost:1234
	go func() {
		fmt.Printf("Start '%s' to listen on inputs\n", spec.Name)
		err = http.ListenAndServe(createAddress(spec), nil)
		if err != nil {
			panic(err)
		}
	}()

	go func() {

		send := func(name string, conn *websocket.Conn) {
			if conn != nil {

				msg := "heyho - " + name
				err := websocket.Message.Send(conn, msg)

				now := time.Now().Format(time.StampMilli)
				if err != nil {
					fmt.Printf("%s - Send a message error: %v\n", now, err)
				} else {
					fmt.Printf("%s - Send a message '%s' to attribute '%s'\n", now, msg, name)
				}
			}

			time.Sleep(200 * time.Millisecond)
		}

		for {

			for k, v := range connectionMap {
				send(k, v)
			}

		}
	}()

	<-done
}

func createOutputConnections(spec Spec) {
	for _, conn := range spec.Connections {

		fmt.Printf("Try to connect to %s\n", getDestURL(conn))
		wsConn, err := websocket.Dial(getDestURL(conn), "", getOrigin(conn))
		if err != nil {
			fmt.Printf("Could not connect to %s\n", getDestURL(conn))
			connectionMap[conn.DestAttrName()] = nil
		} else {
			fmt.Printf("Connection to %s established\n", getDestURL(conn))
			// Add connection to map
			// e.g. connectionMap["A"] = ws://127.0.0.1:1234/C
			// Attribute "A" is connect to cell "ws://127.0.0.1:1234" at attribut "C"
			connectionMap[conn.DestAttrName()] = wsConn
		}
	}
}

func isConnected(attrName string) bool {
	if _, ok := connectionMap[attrName]; ok {
		return true
	}
	return false
}

func getDestURL(conn Connection) string {
	url := fmt.Sprintf("%s/%s", conn.DestAddress(), conn.DestAttrName())
	return url
}
func getOrigin(conn Connection) string {
	url := fmt.Sprintf("http://%s/", conn.DestAddress())
	return url
}
