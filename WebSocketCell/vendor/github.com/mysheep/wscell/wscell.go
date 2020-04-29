package wscell

import (
	"fmt"
	"io"
	"log"
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
	// e.g. can be connect via
	//  A.."ws://localhost:1234/inputA"
	//  B.."ws://localhost:1234/inputB"
	//  C.."ws://localhost:1234/inputC"
	//  D.."ws://localhost:1234/inputD"
	createInputs(spec)

	// Create OUTPUT connections
	createOutputConnections(spec)

	// CleanUp at end
	defer closeConnections(connectionMap)

	// Listen to clients connecting
	// ws://{ip}:{port}/{Name}
	// e.g.
	//   ws://localhost:1234
	go func() {
		log.Printf("Start '%s' to listen on inputs\n", spec.Name)
		err = http.ListenAndServe(createAddress(spec), nil)
		if err != nil {
			panic(err)
		}
	}()

	go func() {

		printDebug := func(msg string, err error, name string) {
			if err != nil {
				log.Printf("Send a message error: %v\n", err)
			} else {
				log.Printf("Send a message '%s' to attribute '%s'\n", msg, name)
			}
		}

		sendTest := func(name string, conn *websocket.Conn) {
			if conn != nil {
				msg := "heyho - " + name
				err := websocket.Message.Send(conn, msg)
				printDebug(msg, err, name)
			}
		}

		time.Sleep(1 * time.Second)

		for i := 0; i < 100000; i++ {
			// TODO: MultiCast
			for attrName, conn := range connectionMap { // Sent to all
				sendTest(attrName, conn)
			}
			//time.Sleep( * time.Millisecond)
		}

	}()

	<-done
}

func createInputHandler(inputName string) func(*websocket.Conn) {

	printMessage := func(message string, err error) {
		if err != nil {
			log.Printf("Received an error '%v'\n", err)
		} else {
			log.Printf("Received message '%s' at attribute '%s'\n", message, inputName)
		}
	}

	inputHandler := func(conn *websocket.Conn) {

		log.Printf("Start inputhandler for attribute '%s'\n", inputName)

		message := ""
		for {
			// Wait
			err := websocket.Message.Receive(conn, &message)
			// EOF -> leave -> Client went away
			if err == io.EOF {
				break
			}
			// Print
			printMessage(message, err)
		}

		log.Printf("Leave inputhandler for attribute '%s'\n", inputName)
	}

	return inputHandler
}

func createInputs(spec Spec) {

	printDebug := func(name string) {
		log.Printf("Create input connection point at 'ws://%s:%s/%s'\n", spec.IP, spec.Port, name)
	}

	for _, attr := range spec.Attributes {
		printDebug(attr.Name())
		http.Handle("/"+attr.Name(), websocket.Handler(createInputHandler(attr.Name())))
	}
}

func createAddress(spec Spec) string {
	address := spec.IP + ":" + spec.Port
	return address
}

func createOutputConnections(spec Spec) {
	for _, conn := range spec.Connections {

		log.Printf("Try to connect to %s\n", getDestURL(conn))
		wsConn, err := websocket.Dial(getDestURL(conn), "", getOrigin(conn))
		if err != nil {
			log.Printf("Could not connect to %s\n", getDestURL(conn))
			connectionMap[conn.DestAttrName()] = nil
		} else {
			log.Printf("Connection to %s established\n", getDestURL(conn))
			// Add connection to map
			// e.g. connectionMap["A"] = ws://127.0.0.1:1234/C
			// Attribute "A" is connect to cell "ws://127.0.0.1:1234" at attribut "C"
			connectionMap[conn.DestAttrName()] = wsConn
		}
	}
}

func closeConnections(conns map[string]*websocket.Conn) {
	for _, conn := range conns {
		if conn != nil {
			log.Println("Close connection")
			conn.Close()
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
