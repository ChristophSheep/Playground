package main

import (
	"flag"
	"net/url"

	"github.com/mysheep/attribute"
	"github.com/mysheep/wscell"
)

// Empty
const EMPTY = "empty"

// Parameter of cell
type Parameter struct {
	Name   string
	Port   string
	PortTo string
}

func getParameter() Parameter {

	portPtr := flag.String("port", "7777", "the port")
	namePtr := flag.String("name", "Adder", "the name")
	portToPtr := flag.String("portto", EMPTY, "the port connect to")

	flag.Parse()

	return Parameter{
		Name:   *namePtr,
		Port:   *portPtr,
		PortTo: *portToPtr,
	}
}

func createURL(scheme string, host string, port string, path string) string {

	if host == "" {
		host = "localhost"
	}

	if scheme == "" {
		scheme = "ws"
	}

	url := url.URL{
		Scheme: scheme,
		Host:   host + ":" + port,
		Path:   path,
	}

	return url.String()
}

func createConnURL(port string) string {
	return createURL("", "", port, "")
}

//
//           *CELL1*					   *CELL2*
//       +------------+					+------------+
//	  -->o /inA  outB o---->  ....  --->o /inA  outA o---->
//	     |            | 				|            |
//	  -->o /inB  outB o---->  ....  --->o /inB  outB o---->
//       +------------+					+------------+
//         [IP]:[PORT]		webSocket	  [IP]:[PORT]
//
//    server        client			 server        client
//
func main() {

	// Switch logging off
	//
	//log.SetOutput(ioutil.Discard)

	// TODO:
	// - So we need a JSON Config:
	// 		- Config has information about cell
	// 		- Config has information about connecting to

	// Command line params
	//
	params := getParameter()

	// Create attributes
	//
	names := []string{"A", "B", "C", "D"}

	attrs := []attribute.Attribute{}
	for _, name := range names {
		attr := attribute.CreateIntAttribute(name)
		attrs = append(attrs, attr)
	}

	// Create connections
	//

	// --------------------
	// 	+---+
	// 	| A |
	// 	| B |      +---+
	// 	| C |o---->| A |o
	// 	| D |o---->| B |o
	// 	+---+	   | C |o
	// 	Cell1	   | D |o
	//             +---+
	//  		   Cell2
	// --------------------

	conns := []wscell.Connection{}
	if params.PortTo != EMPTY {

		connA := wscell.CreateConnection("C", createConnURL(params.PortTo), "A")
		connB := wscell.CreateConnection("B", createConnURL(params.PortTo), "B")
		conns = []wscell.Connection{
			connA,
			connB,
		}
	}

	//  cell://"adder":localhost:7777/{A,B}|{C->A, D->B}

	// TODO:
	// - CalcFns -> out = fn(in)
	//
	spec := wscell.Spec{
		IP:          "localhost", // TODO: always localhost or??
		Port:        params.Port,
		Name:        params.PortTo,
		Attributes:  attrs,
		Connections: conns,
	}

	// Create cell and listen
	//
	wscell.CreateAndListen(spec)

}
