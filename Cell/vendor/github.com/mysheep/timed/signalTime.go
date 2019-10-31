package timed

import (
	"fmt"
	"time"
)

// SignalTime is a boolean signal value with a time
type SignalTime struct {
	val  bool
	time time.Time
}

// MakeSignalTime create a SignalTime object
func MakeSignalTime(val bool, time time.Time) SignalTime {
	return SignalTime{val: val, time: time}
}

// Val get the value of a timed signal value
func (c *SignalTime) Val() bool {
	return c.val
}

// Time get the time of a timed signal value
func (c *SignalTime) Time() time.Time {
	return c.time
}

// String implements the Stringer interface
func (c *SignalTime) String() string {
	return fmt.Sprintf("{val: %t, time: %s}", c.val, c.time.Format(TIME_FORMAT))
}
