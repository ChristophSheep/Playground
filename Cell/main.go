package main

import (
	"github.com/mysheep/cell"
	"github.com/mysheep/cell/integer"
)

func main() {

	//
	//  created with http://asciiflow.com/
	//
	//  in           out
	//      +------+     +---------+
	// ---->| Add1 o---->| Display |
	//      +------+     +---------+

	//          +-------+     +---------+
	// inX ---->|       | out |         |
	//          |  Add  o---->| Display |
	// inY ---->|       | out |         |
	//          +-------+     +---------+

	inX := make(chan int)
	inY := make(chan int)
	out := make(chan int)

	done := make(chan bool)
	waitUntilDone := func() { <-done }

	//
	// Console Commands
	//
	cmds := map[string]func(){
		"quit": func() { done <- true },
		"emit": func() { inX <- 1; inY <- 2 },
	}

	//
	// Setup Network
	//

	addXY := func(x, y int) int { return x + y }

	// go integer.Add1(in, out)
	// OR
	go integer.Lambda2(inX, inY, out, addXY)
	go integer.Display(out)
	go cell.Console(cmds)

	// Wait until Done
	//
	waitUntilDone()
	//
	// Wait until Done
}
