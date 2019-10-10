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
	//M := 10
	ins := make([]chan int, N)
	out2 := make(chan int, 10)
	agg := make(chan int, 10)
	//out := make(chan int)
	for i := 0; i < N; i++ {
		ins[i] = make(chan int)
	}
	go integer.Aggregate(ins, agg, out2)
	go integer.Display(out2)

	//
	// Console Commands
	//
	cmds := map[string]func(){
		"quit": func() { done <- true },
		"emit": func() { inX <- 1; inY <- 2 },
		"emit2": func() {
			for i := 0; i < N; i++ {
				fmt.Println("send", i)
				ins[i] <- i
			}
		},
		"emit3": func() {
			ins = append(ins, make(chan int))
			M := len(ins)

			// TODO: Update - if add a new ins channel, you need also
			go func(i int, ch chan int) {
				for val := range ch {
					agg <- val
				}
			}(M-1, ins[M-1])
			// TODO: Update

			fmt.Println("M:", M)
			for i := 0; i < M; i++ {
				fmt.Println("send", i)
				ins[i] <- i
			}

			// TODO: Update Agg goroutines
		},
	}

	go cell.Console(cmds)

	// Wait until Done
	//
	waitUntilDone()
	//
	// Wait until Done
}
