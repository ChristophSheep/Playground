package main

import (
	"math/rand"

	"github.com/mysheep/cell"
	"github.com/mysheep/cell/integer"
)

func main() {

	//
	//  created with http://asciiflow.com/
	//
	//  	     +--------+
	//  in0	+---->        | out0
	//  		 |   &    o----+
	//  in1	+---->        |    |    +--------+
	//  		 +--------+    +---->        | out2
	//  							|   &    o--->
	//  		 +--------+    +---->        |
	//  in3	+---->        |    |    +--------+
	//  		 |   &    o----+
	//  in4	+---->        | out1
	//  		 +--------+
	//

	// Create input channels
	//
	const INS = 4
	const BUFFER_SIZE = 10
	in := make([]chan int, INS)

	for i := 0; i < INS; i++ {
		in[i] = make(chan int, BUFFER_SIZE)
	}

	// Create output channels
	//
	const OUTS = 3
	out := make([]chan int, OUTS)

	for i := 0; i < OUTS; i++ {
		out[i] = make(chan int, BUFFER_SIZE)
	}

	// Create done channel and waiter
	done := make(chan bool)
	waitUntilDone := func() { <-done }

	//
	// Console Commands Functions
	//
	cmds := map[string]func(){

		"quit": func() { done <- true },

		"emit": func() {
			for i := 0; i < 5; i++ {
				n := rand.Int31n(INS) // [0, n)
				in[n] <- rand.Intn(100)
			}
		},

		"emit1": func() {
			in[0] <- 1
		},

		"emit2": func() {
			in[1] <- 2
		},
	}

	//
	// Setup Network
	//

	go integer.Add2Async(in[0], in[1], out[0])
	go integer.Add2Async(in[2], in[3], out[1])
	go integer.Add2Async(out[0], out[1], out[2])

	go integer.Display(out[2])
	go cell.Console(cmds)

	waitUntilDone()
}
