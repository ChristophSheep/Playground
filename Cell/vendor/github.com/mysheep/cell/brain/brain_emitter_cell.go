package brain

// ----------------------------------------------------------------------------
// EmitterCell has only outputs
// ----------------------------------------------------------------------------

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
	N := 0

	c := EmitterCell{
		name:    name,
		outputs: make([]chan int, N),
	}

	return &c
}
