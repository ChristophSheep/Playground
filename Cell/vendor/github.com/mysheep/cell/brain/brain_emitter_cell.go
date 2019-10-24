package brain

import (
	"time"

	"github.com/mysheep/timed"
)

// ----------------------------------------------------------------------------
// EmitterCell has only outputs
// ----------------------------------------------------------------------------

// +---------+ outputs
// | Emitter +--------->
// +---------+

// ----------------------------------------------------------------------------
// Public
// ----------------------------------------------------------------------------

type EmitterCell struct {
	name    string
	outputs []chan timed.SignalTime
}

func (c *EmitterCell) Name() string {
	return c.name
}

func (c *EmitterCell) OutputConnect(ch chan timed.SignalTime /*not used*/) {
	c.outputs = append(c.outputs, ch)
}

func (c *EmitterCell) EmitOne(t time.Time) {
	one := timed.MakeSignalTime(true, t)
	for _, ch := range c.outputs {
		ch <- one
	}
}

func MakeEmitterCell(name string) *EmitterCell {
	return &EmitterCell{
		name:    name,
		outputs: make([]chan timed.SignalTime, 0),
	}
}
