package brain

import "fmt"

//
// DisplayCell has only inputs
//

type DisplayCell struct {
	name   string
	inputs []chan int
}

func (c *DisplayCell) Name() string {
	return c.name
}

func (c *DisplayCell) AddInput(ch chan int, weight int /*not used*/) {
	c.inputs = append(c.inputs, ch)
}

func (c *DisplayCell) AddOutput(ch chan int /*not used*/) {
	// Display prints to console
}

func (c *DisplayCell) Update() {
	for j := 0; j < len(c.inputs); j++ {
		go Display(c.inputs[j], fmt.Sprintf("display %d", j))
	}
}

func MakeDisplayCell(name string) *DisplayCell {
	N := 0

	c := DisplayCell{
		name:   name,
		inputs: make([]chan int, N),
	}

	return &c
}
