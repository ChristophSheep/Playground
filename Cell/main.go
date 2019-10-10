package main

import (
	"fmt"

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
	// Setup Network
	//

	addXY := func(x, y int) int { return x + y }
	go integer.Lambda2(inX, inY, out, addXY)
	go integer.Display(out)

	//             +---------+
	// ins[0] ---->|         | out +---------+
	//             |  	     |     |		 |
	// ...         |  Aggr.  o---->| Display |
	//             |  	     |     |		 |
	// ins[n] ---->|         | out +---------+
	//             +---------+

	//
	// Aggregation
	//
	N := 5
	ins := make([]chan int, N)
	fin := make(chan int, 10)
	//agg := make(chan int, 10)
	for i := 0; i < N; i++ {
		ins[i] = make(chan int)
	}
	updateFn, aggFn, quitFn := integer.MakeAgg(&ins, fin)
	updateFn()
	go aggFn()
	go integer.Display(fin)

	var addOneFn = func() { ins = append(ins, make(chan int)) }

	//
	// Console Commands
	//
	cmds := map[string]func(){
		"quit": func() { done <- true },
		"emit": func() { inX <- 1; inY <- 2 },
		"agg": func() {
			for i := 0; i < len(ins); i++ {
				fmt.Println("send", i)
				ins[i] <- i
			}
		},
		"add": func() {
			addOneFn()
			fmt.Println("add", len(ins), "ins")
		},
		"exit": func() {
			quitFn()
		},
		"upd": func() {
			updateFn()
		},
		"all": func() {
			quitFn()
			addOneFn()
			updateFn()
		},
	}

	go cell.Console(cmds)

	// Wait until Done
	//
	waitUntilDone()
	//
	// Wait until Done
}
