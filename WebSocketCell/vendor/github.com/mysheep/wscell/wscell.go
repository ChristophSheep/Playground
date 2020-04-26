package wscell

import (
	"fmt"
	"net/http"

	"github.com/mysheep/attribute"
	"golang.org/x/net/websocket"
)

// Spec TODO
type Spec struct {
	IP         string
	Port       string
	Attributes []attribute.Attribute
}

// Cell TODO
type Cell struct {
	Name       string
	IP         string
	Port       string
	Attributes []attribute.Attribute
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

func createInputs(spec Spec) error {

	// Create an input connection point
	// for each attributes
	// ws://{ip}:{port}/{Name}
	// ws://localhost:1234/A
	for _, attr := range spec.Attributes {
		http.Handle("/"+attr.Name(), websocket.Handler(createInputHandler(attr.Name())))
	}

	// Listen to clients connecting
	err := http.ListenAndServe(createAddress(spec), nil)

	// Show Error
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
	return err
}

func createAddress(spec Spec) string {
	address := spec.IP + ":" + spec.Port
	return address
}

func createOutputs(spec Spec) error {
	// TODO
	return nil
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
func Create(spec Spec) error {

	var err error

	// Create INPUT connection points
	// e.g. ws://localhost:1234/inputA
	err = createInputs(spec)

	// Show Error
	if err != nil {
		return err
	}

	// create OUTPUTS connection points
	err = createOutputs(spec)

	// Show Error
	if err != nil {
		return err
	}

	return nil
}
