package brain

import (
	"fmt"
	"time"
)

// ----------------------------------------------------------------------------
// Multi weighted inputs and one output axon with multi connections cell
//
// see https://en.wikipedia.org/wiki/Multipolar_neuron

//     A multipolar neuron (or multipolar neurone) is a type of neuron that
//     possesses a single axon and many dendrites (and dendritic branches),
//     allowing for the integration of a great deal of information from other
//     neurons. These processes are projections from the nerve cell body.
//     Multipolar neurons constitute the majority of neurons in the central
//     nervous system. They include motor neurons and interneurons/relaying
//     neurons are most commonly found in the cortex of the brain, the spinal
//     cord, and also in the autonomic ganglia.

// ----------------------------------------------------------------------------
//
//    dendrides         synapses			axon
//
//                         +-----------+
// ---input[0]--weight[0]->o	       |          +--> outputs[0]
// ---input[1]--weight[1]->o cell body o----------|--> outputs[1]
//    ...                  |   soma    |          |    ...
// ---input[n]--weight[n]->o           |          +--> outputs[n]
//                         +-----------+
//
//                      weights
//
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

const MAXAGE = 1 // max age in seconds of soma in multicell

// ----------------------------------------------------------------------------
// Public
// ----------------------------------------------------------------------------

type weightInput struct {
	weight float64
	input  chan SignalTime
}

type MultiCell struct {
	name string

	// TODO: INputs and Weights should be a struct
	//inputs  []chan SignalTime
	//weights []float64

	weightedInputs []weightInput

	outputs []chan SignalTime

	bodyIn  chan FloatTime  // Drendrid, collect all weighted inputs
	bodyOut chan SignalTime // Fire out to axon
}

func (c *MultiCell) SetWeight(i int, weight float64) {
	if i >= 0 && i < len(c.weightedInputs) {
		c.weightedInputs[i].weight = weight
	}
}

func (c *MultiCell) Weights() []float64 {
	weights := make([]float64, len(c.weightedInputs))
	for i := 0; i < len(c.weightedInputs); i++ {
		weights[i] = c.weightedInputs[i].weight
	}
	return weights
}

func (c *MultiCell) Name() string {
	return c.name
}

func (c *MultiCell) OutputConnect(ch chan SignalTime) {
	c.outputs = append(c.outputs, ch)
}

func (c *MultiCell) InputConnect(ch chan SignalTime, weight float64) {
	// TODO: Input and Weights should be a struct
	c.weightedInputs = append(c.weightedInputs, weightInput{weight: weight, input: ch})
	//c.inputs = append(c.inputs, ch)
	//c.weights = append(c.weights, weight)
	lastIndex := len(c.weightedInputs) - 1
	go c.synapse(lastIndex)
}

/*
	Make a multi input multi output cell
*/
func MakeMultiCell(name string, threshold float64) *MultiCell {

	c := MultiCell{
		name: name,

		//inputs:  make([]chan SignalTime, 0, 1024),
		//weights: make([]float64, 0, 1024),

		weightedInputs: make([]weightInput, 0, 1024),
		outputs:        make([]chan SignalTime, 0, 1024),

		bodyIn:  make(chan FloatTime, 1024), // size of inputs, so that they not need to wait
		bodyOut: make(chan SignalTime, 10),
	}

	go c.soma(threshold)
	go c.axon()

	return &c
}

// ----------------------------------------------------------------------------
// Private
// ----------------------------------------------------------------------------

/*
	Weighted synapae
*/
func (c *MultiCell) synapse(index int) func() {

	for {
		signal := <-c.weightedInputs[index].input
		val := FloatTime{val: 0.0, time: signal.time}
		if signal.val {
			val = FloatTime{val: c.weightedInputs[index].weight, time: signal.time}
		}
		c.bodyIn <- val
	}
}

/*
	Soma or also called cell body with threshold to fire
*/
func (c *MultiCell) soma(threshold float64) {

	const EPSILON = 0.5 // TODO Float64 has problems

	var (
		fire = func() { c.bodyOut <- SignalTime{val: true, time: time.Now()} }
		sums = MakeFloatSums(MAXAGE)
	)

	for {
		select {
		case val := <-c.bodyIn:
			sums.AddVal(val)

			if sum, ok := sums.getSum(val.time); ok && sum >= (threshold-EPSILON) {
				fmt.Printf("%s -> sum: %f, threshold: %f", c.name, sum, threshold)
				fire()                  // fire action potential
				sums.resetSum(val.time) // go back to rest level
			}
		}
	}
}

/*
	Axon of a multi cell
*/
func (c *MultiCell) axon() {
	for {
		val := <-c.bodyOut
		for _, out := range c.outputs {
			// Fire to all outputs with goroutines
			// to rush out the signal
			go func(out chan SignalTime, val SignalTime) {
				out <- val
			}(out, val)
		}
	}
}
