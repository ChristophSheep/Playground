package brain

import (
	"fmt"
)

//         Synapses
//  inputs +->|	          | outputs   +-->
// --------+->| cell body |-----------+-->
//  inputs +->|           |  Axon     +-->
//         Weights

type Cell struct {

	// public
	Name    string
	Inputs  []chan int
	Weights []int
	Outputs []chan int

	// private
	bodyIn      chan int
	cellOutAxIn chan int
}

func MakeCell(name string) *Cell {

	N := 0 // No input per default for now

	c := Cell{
		Name:        name,
		Inputs:      make([]chan int, N),
		Weights:     make([]int, N),
		Outputs:     make([]chan int, N),
		bodyIn:      make(chan int),
		cellOutAxIn: make(chan int),
	}

	for j := 0; j < len(c.Inputs); j++ {
		go Synapse(c.Weights[j], c.Inputs[j], c.bodyIn)
	}

	go Body(c.bodyIn, c.cellOutAxIn)
	go Axon2(c.cellOutAxIn, &c.Outputs)

	return &c
}

func CellByName(name string) *Cell {
	return &Cell{}
}

func (c *Cell) update() {
	for j := 0; j < len(c.Inputs); j++ {
		go Synapse(c.Weights[j], c.Inputs[j], c.bodyIn)
	}
}

func (c *Cell) AddOutput(ch chan int) {
	c.addOutput(ch)
}

func (c *Cell) addOutput(ch chan int) {
	c.Outputs = append(c.Outputs, ch)
}

func (c *Cell) AddInput(ch chan int, weight int) {
	c.addInput(ch, weight)
}

func (c *Cell) addInput(ch chan int, weight int) {
	c.Inputs = append(c.Inputs, ch)
	c.Weights = append(c.Weights, weight)
}

func (cFr *Cell) ConnectWith(cTo *Cell, weight int) {
	//
	//  outputs       inputs
	//  (from) -----> (to)
	//           ch    weights

	ch := make(chan int)

	cFr.addOutput(ch)
	cTo.addInput(ch, weight)

	cFr.update()
	cTo.update()

	fmt.Println("cell", cFr.Name, "with cell", cTo.Name, "connected")
}
