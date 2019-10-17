package brain

import "time"

// ----------------------------------------------------------------------------
// EmitterCell has only outputs
// ----------------------------------------------------------------------------

// +---------+ outputs
// | Emitter +--------->
// +---------+

type EmitterCell struct {
	name    string
	outputs []chan IntTime
}

func (c *EmitterCell) Name() string {
	return c.name
}

func (c *EmitterCell) OutputConnect(ch chan IntTime /*not used*/) {
	c.outputs = append(c.outputs, ch)
}

func (c *EmitterCell) EmitOne(t time.Time) {
	for _, out := range c.outputs {
		out <- IntTime{val: 1, time: t}
	}
}

func MakeEmitterCell(name string) *EmitterCell {
	return &EmitterCell{
		name:    name,
		outputs: make([]chan IntTime, 0),
	}
}
