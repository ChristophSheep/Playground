package main

import (
	"flag"

	"github.com/mysheep/attribute"
	"github.com/mysheep/wscell"
)

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
func main() {

	// TODO:
	// - So we need a JSON Config
	// - Config has information about cell
	// - Config has information about connect to

	// Parse parameter
	portPtr := flag.String("port", "1234", "the port")
	namePtr := flag.String("name", "Adder", "the name")
	portToPtr := flag.String("portto", "", "the port connect to")
	flag.Parse()

	// Create input Attributes
	//
	attrA := attribute.CreateIntAttribute("A")
	attrB := attribute.CreateIntAttribute("B")
	attrC := attribute.CreateIntAttribute("C") // c = a + b
	attrD := attribute.CreateIntAttribute("D") // d = a - b

	portTo := *portToPtr
	conns := []wscell.Connection{}

	if portTo != "" {
		connA := wscell.CreateConnection("C", "ws://127.0.0.1:"+portTo, "A")
		connB := wscell.CreateConnection("B", "ws://127.0.0.1:"+portTo, "B")
		conns = []wscell.Connection{
			connA,
			connB,
		}
	}

	// TODO:
	// - CalcFns -> out = fn(in)
	//
	spec := wscell.Spec{
		IP:   "localhost",
		Port: *portPtr,
		Name: *namePtr,
		Attributes: []attribute.Attribute{
			attrA,
			attrB,
			attrC,
			attrD,
		},
		Connections: conns,
	}

	// Create cell and listen
	//
	wscell.CreateAndListen(spec)

}
