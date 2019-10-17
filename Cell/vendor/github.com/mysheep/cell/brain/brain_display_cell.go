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
	inputs []chan IntTime
}

func (c *DisplayCell) Name() string {
	return c.name
}

func (c *DisplayCell) InputConnect(ch chan IntTime, weight float64 /*not used*/) {
	c.inputs = append(c.inputs, ch)
	go Display(ch, fmt.Sprintf("Cell '%s' has fired", c.Name()))
}

func MakeDisplayCell(name string) *DisplayCell {
	return &DisplayCell{
		name:   name,
		inputs: make([]chan IntTime, 0),
	}
}
