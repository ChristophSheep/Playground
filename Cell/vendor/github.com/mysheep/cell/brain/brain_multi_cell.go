package brain

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

type MultiCell struct {
	name string

	inputs  []chan IntTime
	weights []float64
	outputs []chan IntTime

	bodyIn  chan FloatTime
	bodyOut chan IntTime
}

func MakeMultiCell(name string, threshold float64) *MultiCell {

	c := MultiCell{
		name: name,

		inputs:  make([]chan IntTime, 0),
		weights: make([]float64, 0),
		outputs: make([]chan IntTime, 0),

		bodyIn:  make(chan FloatTime, 100), // buffered, because many pipe in
		bodyOut: make(chan IntTime),
	}

	go soma(&c, threshold)
	go axon(&c)

	return &c
}

func soma(c *MultiCell, threshold float64) {

	const MAXAGE = 1 // seconds
	sums := MakeFloatSums(1)

	/*
		var sendOut = func() {

			var fireUntil = func() {
				for ; sum >= threshold; sum = sum - threshold {
					c.bodyOut <- IntTime{val: 1, time: time.Now()}
				}
			}

			var rest = func() {
				//c.bodyOut <- 0 // TODO: Rethink - fire a ZERO ??
			}

			if sum > threshold {
				fireUntil()
			} else {
				rest()
			}
		}
	*/

	for {
		select {
		case val := <-c.bodyIn:
			sums.AddVal(val)
			// TODO: sendOut()
		}
	}
}

func axon(c *MultiCell) {
	for {
		val := <-c.bodyOut
		for _, out := range c.outputs {
			out <- val
		}
	}
}

func (c *MultiCell) SetWeight(i int, weight float64) {
	c.weights[i] = weight
}

func (c *MultiCell) Weights() []float64 {
	return c.weights
}

func (c *MultiCell) Name() string {
	return c.name
}

func (c *MultiCell) OutputConnect(ch chan IntTime) {
	c.outputs = append(c.outputs, ch)
}

func (c *MultiCell) InputConnect(ch chan IntTime, weight float64) {
	c.inputs = append(c.inputs, ch)
	c.weights = append(c.weights, weight)
	last := len(c.weights) - 1
	// TODO: Wheight must be a pointer that it could change over time
	go Synapse(&c.weights[last], ch, c.bodyIn)
}
