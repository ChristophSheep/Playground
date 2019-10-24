package brain

import (
	"fmt"

	"github.com/mysheep/timed"
)

// ----------------------------------------------------------------------------
// DisplayCell has only inputs
// ----------------------------------------------------------------------------

//  inputs  +---------+
// -------->| Display |
//          +---------+

// ----------------------------------------------------------------------------
// Public
// ----------------------------------------------------------------------------

type DisplayCell struct {
	name   string
	inputs []chan timed.SignalTime
}

func MakeDisplayCell(name string) *DisplayCell {
	return &DisplayCell{
		name:   name,
		inputs: make([]chan timed.SignalTime, 0),
	}
}

func (c *DisplayCell) Name() string {
	return c.name
}

func (c *DisplayCell) InputConnect(ch chan timed.SignalTime, weight float64 /*not used*/) {
	c.inputs = append(c.inputs, ch)
	go display(ch, fmt.Sprintf("Cell '%s' has fired", c.Name()))
}

// ----------------------------------------------------------------------------
// Private
// ----------------------------------------------------------------------------

func display(in <-chan timed.SignalTime, text string) {
	for {
		x := <-in
		fmt.Println(getNow(), "-", text, x.String())
	}
}
