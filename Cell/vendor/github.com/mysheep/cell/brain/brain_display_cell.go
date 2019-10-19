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
	inputs []chan SignalTime
}

func (c *DisplayCell) Name() string {
	return c.name
}

func (c *DisplayCell) InputConnect(ch chan SignalTime, weight float64 /*not used*/) {
	c.inputs = append(c.inputs, ch)
	go Display(ch, fmt.Sprintf("Cell '%s' has fired", c.Name()))
}

func MakeDisplayCell(name string) *DisplayCell {
	return &DisplayCell{
		name:   name,
		inputs: make([]chan SignalTime, 0),
	}
}
