package brain

import (
	"time"

	"github.com/mysheep/timed"
)

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

const TIME_FORMAT = "15:04:05.00000"

// ----------------------------------------------------------------------------
// Interfaces
// ----------------------------------------------------------------------------

type InputConnector interface {
	InputConnect(ch chan timed.SignalTime, weight float64)
}

type OutputConnector interface {
	OutputConnect(ch chan timed.SignalTime)
}

type Namer interface {
	Name() string
}

// ----------------------------------------------------------------------------
// Public
// ----------------------------------------------------------------------------

func ConnectBy(out OutputConnector, in InputConnector, weight float64) {
	ch := make(chan timed.SignalTime, 10)
	out.OutputConnect(ch)
	in.InputConnect(ch, weight)
}

// ----------------------------------------------------------------------------
// Private
// ----------------------------------------------------------------------------

func getNow() string {
	return time.Now().Format(TIME_FORMAT)
}
