package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/mysheep/cell"
	"github.com/mysheep/cell/brain"
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
	// inX ---->|       | res |         |
	//          |  Add  o---->| Display |
	// inY ---->|       | res |         |
	//          +-------+     +---------+

	inX := make(chan int)
	inY := make(chan int)
	res := make(chan int)

	done := make(chan bool)
	waitUntilDone := func() { <-done }

	//
	// Setup Network
	//

	addXY := func(x, y int) int { return x + y }
	go integer.Lambda2(inX, inY, res, addXY)
	go integer.Display(res)

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

	addFn, aggFn := integer.MakeDynAgg(&ins, fin)

	go aggFn()
	go integer.Display(fin)

	var addOneFn = func() {
		newCh := make(chan int)
		addFn(newCh)
	}

	//
	// Cell with weighted Synapses, Body and Axon
	//

	S := 5

	sIns := make([]chan int, S)
	sOuts := make([]chan int, S)
	weights := make([]int, S)

	for j := 0; j < S; j++ {
		sIns[j] = make(chan int)
		sOuts[j] = make(chan int)
		weights[j] = rand.Intn(7)
		go brain.Synapse(weights[j], sIns[j], sOuts[j])
	}

	A := 1
	axIn := make(chan int)
	axOuts := make([]chan int, A)

	for j := 0; j < A; j++ {
		axOuts[j] = make(chan int)
		go brain.Display(axOuts[j], "- axon")
	}
	go brain.Cell(sOuts, axIn)
	go brain.Axon(axIn, axOuts)

	//
	// Console Commands
	//
	cmds := map[string]func(){
		"quit": func() { done <- true },
		"exit": func() { done <- true },
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
		"add10": func() {
			N := 10
			for i := 0; i < N; i++ {
				addOneFn()
			}
		},
		"cell": func() {
			for ii := 0; ii < 100; ii++ {
				i := rand.Intn(S)
				fmt.Println("synapse ", i, "weight:", weights[i])
				sIns[i] <- i
				time.Sleep(50 * time.Millisecond)
			}
		},
	}

	go cell.Console(cmds)

	// Wait until Done
	//
	waitUntilDone()
	//
	// Wait until Done
}
