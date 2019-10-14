package brain

import "fmt"

// ----------------------------------------------------------------------------
// DisplayCell has only inputs
// ----------------------------------------------------------------------------

//  inputs  +---------+
// -------->| Display |
//          +---------+

type DisplayCell struct {
	name   string
	inputs []chan int
}

func (c *DisplayCell) Name() string {
	return c.name
}

func (c *DisplayCell) InputConnect(ch chan int, weight float64 /*not used*/) {
	c.inputs = append(c.inputs, ch)
	go Display(ch, fmt.Sprintf("display '%s' %p", c.Name(), &c))
}

func MakeDisplayCell(name string) *DisplayCell {
	return &DisplayCell{
		name:   name,
		inputs: make([]chan int, 0),
	}
}
