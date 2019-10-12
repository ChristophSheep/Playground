package brain

// ----------------------------------------------------------------------------
// Multi weighted inputs and one output axon with multi connections cell
// ----------------------------------------------------------------------------
//            synapses
//                +-----------+
// ---input[0]--->o	          |         +--> outputs[0]
// ---input[0]--->o cell body o---------|--> outputs[1]
// ---input[n]--->o   soma    |   axon  +--> outputs[n]
//                +-----------+
//             weights
// ----------------------------------------------------------------------------

type Cell struct {
	name string

	inputs  []chan int
	weights []int
	outputs []chan int

	bodyIn chan int
	axIn   chan int
}

func MakeCell(name string) *Cell {

	N := 0 // No input per default for now

	c := Cell{
		name: name,

		inputs:  make([]chan int, N),
		weights: make([]int, N),
		outputs: make([]chan int, N),

		bodyIn: make(chan int, 100), // buffered, because many pipe in
		axIn:   make(chan int),
	}

	for j := 0; j < len(c.inputs); j++ {
		go Synapse(c.weights[j], c.inputs[j], c.bodyIn)
	}

	go Soma(c.bodyIn, c.axIn)
	go Axon2(c.axIn, &c.outputs) // TODO: WHY THIS NEED TO BE A POINTERS ??

	return &c
}

func (c *Cell) Name() string {
	return c.name
}

func (c *Cell) Update() {
	for j := 0; j < len(c.inputs); j++ {
		go Synapse(c.weights[j], c.inputs[j], c.bodyIn)
	}
}

func (c *Cell) AddOutput(ch chan int) {
	c.outputs = append(c.outputs, ch)
}

func (c *Cell) AddInput(ch chan int, weight int) {
	c.inputs = append(c.inputs, ch)
	c.weights = append(c.weights, weight)
}
