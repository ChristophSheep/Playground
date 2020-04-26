package main

import (
	"fmt"

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

	cell1 := wscell.Create(spec1)
	fmt.Println(cell1)

}
