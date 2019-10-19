package brain

import "time"

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
	outputs []chan SignalTime
}

func (c *EmitterCell) Name() string {
	return c.name
}

func (c *EmitterCell) OutputConnect(ch chan SignalTime /*not used*/) {
	c.outputs = append(c.outputs, ch)
}

func (c *EmitterCell) EmitOne(t time.Time) {
	one := SignalTime{val: true, time: t}
	for _, ch := range c.outputs {
		ch <- one
	}
}

func MakeEmitterCell(name string) *EmitterCell {
	return &EmitterCell{
		name:    name,
		outputs: make([]chan SignalTime, 0),
	}
}
