package brain

// ----------------------------------------------------------------------------
// EmitterCell has only outputs
// ----------------------------------------------------------------------------

// +---------+ outputs
// | Emitter +--------->
// +---------+

type EmitterCell struct {
	name    string
	outputs []chan int
}

func (c *EmitterCell) Name() string {
	return c.name
}

func (c *EmitterCell) AddInput(ch chan int, weight int /*not used*/) {
	// Emitter has only outputs
}

func (c *EmitterCell) AddOutput(ch chan int /*not used*/) {
	c.outputs = append(c.outputs, ch)
}

func (c *EmitterCell) Update() {
	// Nothing to update
}

func (c *EmitterCell) EmitOne() {
	for _, out := range c.outputs {
		out <- 1
	}
}

func MakeEmitterCell(name string) *EmitterCell {
	c := EmitterCell{
		name:    name,
		outputs: make([]chan int, 0),
	}

	return &c
}
