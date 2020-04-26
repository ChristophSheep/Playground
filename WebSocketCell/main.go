package main

import (
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

	//
	// Info: Only one cell could be create by one process
	//

	// TODO:
	// - So we need a JSON Config
	// - Config has information about cell
	// - COnfig has information about connect to

	attrA := attribute.CreateIntAttribute("A")
	attrB := attribute.CreateIntAttribute("B")

	spec1 := wscell.Spec{
		IP:   "localhost",
		Port: "1234",
		Attributes: []attribute.Attribute{
			attrA,
			attrB,
		},
	}

	// Cell 1
	//
	cell1 := wscell.Create(spec1)

	// TODO
	spec2 := spec1

	// Cell 2
	//
	cell2 := wscell.Create(spec2)

	// TODO

	// cell1.A -----> cell2.A
	//
	cell1.ConnectTo(cell1.GetAttributByName("A"), cell2, cell2.GetAttributByName("A"))
}
