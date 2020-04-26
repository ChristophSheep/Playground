package wscell

import (
	"fmt"
	"net/http"

	"github.com/mysheep/attribute"
	"golang.org/x/net/websocket"
)

// Spec to create a cell
type Spec struct {
	Name       string
	IP         string
	Port       string
	Attributes []attribute.Attribute
}

// Cell interface
type Cell interface {
	Name() string
	IP() string
	Port() string
	Origin() string
	GetAttributByName(string) attribute.Attribute
	ConnectTo(attribute.Attribute, Cell, attribute.Attribute)
}

type connection struct {
	src      Cell
	srcAttr  attribute.Attribute
	dest     Cell
	destAttr attribute.Attribute
	wsConn   *websocket.Conn // WebSocket Connection
}

type cell struct {
	name        string
	ip          string
	port        string
	attributes  []attribute.Attribute
	connections []connection
}

func createInputHandler(inputName string) func(*websocket.Conn) {
	inputHandler := func(ws *websocket.Conn) {

		for {
			var message string
			err := websocket.Message.Receive(ws, &message)
			if err != nil {
				fmt.Printf("Received error '%v'\n", err)
			} else {
				fmt.Printf("Received message '%s' at input [%s]\n", message, inputName)
			}
		}
	}
	fmt.Printf("Create input \"%s\"\n", inputName)
	return inputHandler
}

func createInputs(spec Spec) {

	// Create an input connection point
	// for each attributes
	// ws://{ip}:{port}/{Name}
	// e.g.
	//   ws://localhost:1234/A
	for _, attr := range spec.Attributes {
		http.Handle("/"+attr.Name(), websocket.Handler(createInputHandler(attr.Name())))
	}
}

func createAddress(spec Spec) string {
	address := spec.IP + ":" + spec.Port
	return address
}

// Create TODO
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
func Create(spec Spec) Cell {

	var err error

	// Create INPUT connection points
	// e.g. ws://localhost:1234/inputA
	createInputs(spec)

	// Listen to clients connecting
	// ws://{ip}:{port}/{Name}
	// e.g.
	//   ws://localhost:1234
	err = http.ListenAndServe(createAddress(spec), nil)
	if err != nil {
		panic(err)
	}

	// Cell itself
	//
	c := cell{
		name:       spec.Name,
		ip:         spec.IP,
		port:       spec.Port,
		attributes: spec.Attributes,
	}

	fmt.Printf("Cell %s\n created", spec.Name)
	return c
}

func (c cell) Name() string {
	return c.name
}

func (c cell) IP() string {
	return c.ip
}

func (c cell) Port() string {
	return c.port
}

// Origin TODO
func (c cell) Origin() string {
	return fmt.Sprintf("http://%s:%s", c.ip, c.port)
}

// GetAttributByName TODO
func (c cell) GetAttributByName(name string) attribute.Attribute {
	for _, attr := range c.attributes {
		if attr.Name() == name {
			return attr
		}
	}
	return attribute.Empty
}

func getDestURL(dest Cell, destAttr attribute.Attribute) string {
	url := fmt.Sprintf("ws://%s:%s/%s", dest.IP(), dest.Port(), destAttr.Name())
	return url
}

// ConnectTo TODO
func (c cell) ConnectTo(srcAttr attribute.Attribute, dest Cell, destAttr attribute.Attribute) {

	// Establish connection
	wsConn, err := websocket.Dial(getDestURL(dest, destAttr), "", dest.Origin())
	if err != nil {
		fmt.Println(err)
		return
	}

	// Add connection
	conn := connection{
		src:      c,
		srcAttr:  srcAttr,
		dest:     dest,
		destAttr: destAttr,
		wsConn:   wsConn,
	}
	c.connections = append(c.connections, conn)

}
