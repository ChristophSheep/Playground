package brain

import (
	"fmt"
	"time"
)

type SignalTime struct {
	val  bool
	time time.Time
}

func MakeSignalTime(val bool, time time.Time) SignalTime {
	return SignalTime{val: val, time: time}
}

func (c *SignalTime) String() string {
	return fmt.Sprintf("{val: %t, time: %s}", c.val, c.time.Format(TIME_FORMAT))
}
