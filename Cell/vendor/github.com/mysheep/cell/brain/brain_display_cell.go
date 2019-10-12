package brain

import "fmt"

// ----------------------------------------------------------------------------
// DisplayCell has only inputs
// ----------------------------------------------------------------------------

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
	for _, in := range c.inputs {
		//fmt.Println("display cell update i:", i)
		go Display(in, fmt.Sprintf("display '%s'", c.Name()))
	}
}

func MakeDisplayCell(name string) *DisplayCell {
	c := DisplayCell{
		name:   name,
		inputs: make([]chan int, 0),
	}

	return &c
}
