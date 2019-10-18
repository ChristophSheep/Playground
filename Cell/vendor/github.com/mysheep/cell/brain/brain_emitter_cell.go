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
	one := IntTime{val: 1, time: t}
	for _, ch := range c.outputs {
		ch <- one
	}
}

func (c *EmitterCell) EmitZero(t time.Time) {
	one := IntTime{val: 0, time: t}
	for _, ch := range c.outputs {
		ch <- one
	}
}

func MakeEmitterCell(name string) *EmitterCell {
	return &EmitterCell{
		name:    name,
		outputs: make([]chan IntTime, 0),
	}
}
