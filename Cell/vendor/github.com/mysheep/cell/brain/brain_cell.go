package brain

import (
	"fmt"
)

type Connect interface {
	AddInput(ch chan int, weight int)
	AddOutput(ch chan int)
	Update()
	Name() string
}

//         Synapses
//  inputs +->|	          | outputs   +-->
// --------+->| cell body |-----------+-->
//  inputs +->|           |  Axon     +-->
//         Weights

type Cell struct {

	// public
	name string

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
		name:        name,
		Inputs:      make([]chan int, N),
		Weights:     make([]int, N),
		Outputs:     make([]chan int, N),
		bodyIn:      make(chan int),
		cellOutAxIn: make(chan int),
	}

	for j := 0; j < len(c.Inputs); j++ {
		go Synapse(c.Weights[j], c.Inputs[j], c.bodyIn)
	}

	go Soma(c.bodyIn, c.cellOutAxIn)
	go Axon2(c.cellOutAxIn, &c.Outputs)

	return &c
}

func (c *Cell) Name() string {
	return c.name
}

func (c *Cell) Update() {
	for j := 0; j < len(c.Inputs); j++ {
		go Synapse(c.Weights[j], c.Inputs[j], c.bodyIn)
	}
}

func (c *Cell) AddOutput(ch chan int) {
	c.Outputs = append(c.Outputs, ch)
}

func (c *Cell) AddInput(ch chan int, weight int) {
	c.Inputs = append(c.Inputs, ch)
	c.Weights = append(c.Weights, weight)

}

/*
func (cFr *Cell) ConnectWith(cTo *Cell, weight int) {
	//
	//  outputs       inputs
	//  (from) -----> (to)
	//           ch    weights

	ch := make(chan int)

	cFr.addOutput(ch)
	cTo.addInput(ch, weight)

	cFr.update() // TODO: kill old funcs
	cTo.update() // TODO: kill old funcs

	fmt.Println("cell", cFr.Name(), "with cell", cTo.Name(), "connected")
}
*/
//
// DisplayCell has only inputs
//

type DisplayCell struct {
	name   string
	Inputs []chan int
}

func (c *DisplayCell) Name() string {
	return c.name
}

func (c *DisplayCell) AddInput(ch chan int, weight int /*not used*/) {
	c.Inputs = append(c.Inputs, ch)
	fmt.Println("display cell len(inputs)", len(c.Inputs))
}

func (c *DisplayCell) AddOutput(ch chan int /*not used*/) {
	// Display prints to console
}

func (c *DisplayCell) Update() {
	for j := 0; j < len(c.Inputs); j++ {
		go Display(c.Inputs[j], fmt.Sprintf("display %d", j))
	}
}

func MakeDisplayCell(name string) *DisplayCell {
	N := 0

	c := DisplayCell{
		name:   name,
		Inputs: make([]chan int, N),
	}

	return &c
}

func ConnectBy(from, to Connect, weight int) {
	ch := make(chan int)

	from.AddOutput(ch)
	to.AddInput(ch, weight)

	from.Update()
	to.Update()

	fmt.Println("cell", from.Name(), "with cell", to.Name(), "connected")

}

//
// EmitterCell has only outputs
//

type EmitterCell struct {
	name    string
	Outputs []chan int
}

func (c *EmitterCell) Name() string {
	return c.name
}

func (c *EmitterCell) AddInput(ch chan int, weight int /*not used*/) {
	// Emitter has only outputs
}

func (c *EmitterCell) AddOutput(ch chan int /*not used*/) {
	c.Outputs = append(c.Outputs, ch)
}

func (c *EmitterCell) Update() {
	// Nothing to update
}

func (c *EmitterCell) EmitOne() {
	for _, out := range c.Outputs {
		out <- 1
	}
}

func MakeEmitterCell(name string) *EmitterCell {
	N := 0

	c := EmitterCell{
		name:    name,
		Outputs: make([]chan int, N),
	}

	return &c
}
