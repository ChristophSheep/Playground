package example1

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/mysheep/cell/brain"
	"github.com/mysheep/cell/integer"
	"github.com/mysheep/console"
	"github.com/mysheep/timed"
)

func print(ys []int) {
	for _, y := range ys {
		fmt.Print("y:", y)
	}
	fmt.Println()
}

func Run() {

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

	//     dendride
	//   	          +----------+
	//   0 +------->w1|          |
	//   	          |          |
	//   1 +------->w2|   cell   |  axon
	//   	          |  (soma)  o------->
	//         ...    |   body   |
	//   	          |          |
	//   n +------->wn|          |
	//   	          +----------+
	//            synapses

	S := 5

	bIn := make(chan timed.FloatTime, 100) // buffered body input for aggretion of all synapses
	sIns := make([]chan timed.SignalTime, S)
	weights := make([]float64, S)

	for j := 0; j < S; j++ {
		sIns[j] = make(chan timed.SignalTime)
		weights[j] = float64(rand.Intn(7))
		go brain.Synapse(&weights[j], sIns[j], bIn)
	}

	A := 1
	axIn := make(chan timed.SignalTime)
	axOuts := make([]chan timed.SignalTime, A)

	for j := 0; j < A; j++ {
		axOuts[j] = make(chan timed.SignalTime)
		go brain.Writer(axOuts[j], fmt.Sprintf("out%d", j))
	}
	go brain.Body(bIn, axIn, brain.THRESHOLD)
	go brain.Axon(axIn, axOuts)

	//
	// Create two cells and connect them
	//

	fmt.Println("Setup network: 1 emitter + 2 multi cells + 1 display")

	cell1 := brain.MakeMultiCell("cell_1", 10)
	cell2 := brain.MakeMultiCell("cell_2", 10)

	//  13
	// -->(cell1)      (cell2)--->

	brain.ConnectBy(cell1, cell2, 7)

	display1 := brain.MakeDisplayCell("display_1")
	brain.ConnectBy(cell2, display1, 1)

	emitter1 := brain.MakeEmitterCell("emitter_1")
	brain.ConnectBy(emitter1, cell1, 13)

	//  13        7
	// -->(cell1)----->(cell2)--->Display

	//
	// Example with 3 cells from book Manfred Spitzer
	//

	//   +--------+                   +--------+
	//   |        o------------A--( 5)>        |    +-----------+
	//   | Emit_A o---+    +---B--( 3)> Cell_A o----> Display_C |
	//   |        o-+ |    | +-C--(~3)>        |    +-----------+
	//   +--------+ | |    | |        +--------+
	//
	//   +--------+   |    |         +--------+
	//   |        o---+    +---A--( 5)>        |    +-----------+
	//   | Emit_B o------------B--( 3)> Cell_B o----> Display_B |
	//   |        o---+    +---C--(10)>        |    +-----------+
	//   +--------+   |    |         +--------+
	//
	//   +--------+ | |    | |        +--------+
	//   |        o-+ |    | +-A--( 5)>        |    +-----------+
	//   | Emit_C o---+    +---B--( 3)> Cell_C o----> Display_C |
	//   |        o------------C--(~3)>        |    +-----------+
	//   +--------+                   +--------+

	emitterA := brain.MakeEmitterCell("emitter_A")
	emitterB := brain.MakeEmitterCell("emitter_B")
	emitterC := brain.MakeEmitterCell("emitter_C")

	cellA := brain.MakeMultiCell("cell_A", 8)
	cellB := brain.MakeMultiCell("cell_B", 8)
	cellC := brain.MakeMultiCell("cell_C", 8)

	displayA := brain.MakeDisplayCell("display_A")
	displayB := brain.MakeDisplayCell("display_B")
	displayC := brain.MakeDisplayCell("display_C")

	brain.ConnectBy(emitterA, cellA, 5)
	brain.ConnectBy(emitterA, cellB, 3)
	brain.ConnectBy(emitterA, cellC, -3)

	brain.ConnectBy(emitterB, cellA, -5)
	brain.ConnectBy(emitterB, cellB, 3)
	brain.ConnectBy(emitterB, cellC, 10)

	brain.ConnectBy(emitterC, cellA, 5)
	brain.ConnectBy(emitterC, cellB, 3)
	brain.ConnectBy(emitterC, cellC, -3)

	brain.ConnectBy(cellA, displayA, 1)
	brain.ConnectBy(cellB, displayB, 1)
	brain.ConnectBy(cellC, displayC, 1)

	//
	// Console Commands
	//
	cmds := map[string]func([]string){
		"quit": func(params []string) { done <- true },
		"exit": func(params []string) { done <- true },
		"q":    func(params []string) { done <- true },
		"emit": func(params []string) { inX <- 1; inY <- 2 },
		"agg": func(params []string) {
			for i := 0; i < len(ins); i++ {
				fmt.Println("send", i)
				ins[i] <- i
			}
		},
		"add": func(params []string) {
			addOneFn()
			fmt.Println("add", len(ins), "ins")
		},
		"add10": func(params []string) {
			N := 10
			for i := 0; i < N; i++ {
				addOneFn()
			}
		},
		"cell": func(params []string) {
			one := timed.MakeSignalTime(true, time.Now())
			for ii := 0; ii < 100; ii++ {
				i := rand.Intn(S)
				sIns[i] <- one
				time.Sleep(50 * time.Millisecond)
			}
		},
		"con": func(params []string) {
			now := time.Now()
			for k := 0; k < 10; k++ {
				emitter1.EmitOne(now)
				time.Sleep(50 * time.Millisecond)
			}
		},
		"ex1": func(params []string) {
			now := time.Now()
			emitterA.EmitOne(now)
			emitterB.EmitOne(now)
			emitterC.EmitOne(now)
		},
	}

	go console.Go(cmds)

	// Wait until Done
	//
	waitUntilDone()
	//
	// Wait until Done

	fmt.Println("BYE")
}
