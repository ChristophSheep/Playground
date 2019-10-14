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
	weights []float64
	outputs []chan int

	bodyIn  chan float64
	bodyOut chan int
}

func MakeMultiCell(name string, threshold int) *Cell {

	c := Cell{
		name: name,

		inputs:  make([]chan int, 0),
		weights: make([]float64, 0),
		outputs: make([]chan int, 0),

		bodyIn:  make(chan float64, 100), // buffered, because many pipe in
		bodyOut: make(chan int),
	}

	go soma(&c, threshold)
	go axon(&c)

	return &c
}

func soma(c *Cell, threshold int) {

	sum := 0.0

	var sendOut = func() {

		var fireUntil = func() {
			for ; sum >= float64(threshold); sum = sum - float64(threshold) {
				c.bodyOut <- 1
			}
		}

		var rest = func() {
			//c.bodyOut <- 0 // TODO: Rethink - fire a ZERO ??
		}

		if sum > float64(threshold) {
			fireUntil()
		} else {
			rest()
		}
	}

	for {
		select {
		case val := <-c.bodyIn:
			sum = sum + val
			sendOut()
		}
	}
}

func axon(c *Cell) {
	for {
		val := <-c.bodyOut
		for _, out := range c.outputs {
			out <- val
		}
	}
}

func (c *Cell) SetWeight(i int, weight float64) {
	c.weights[i] = weight
}

func (c *Cell) Weights() []float64 {
	return c.weights
}

func (c *Cell) Name() string {
	return c.name
}

func (c *Cell) OutputConnect(ch chan int) {
	c.outputs = append(c.outputs, ch)
}

func (c *Cell) InputConnect(ch chan int, weight float64) {
	c.inputs = append(c.inputs, ch)
	c.weights = append(c.weights, weight)
	go Synapse(weight, ch, c.bodyIn)
}
