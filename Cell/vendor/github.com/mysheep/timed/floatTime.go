package timed

import (
	"fmt"
	"time"
)

// FloatTime is a float64 value with a time
type FloatTime struct {
	val  float64
	time time.Time
}

// MakeFloatTime create a FloatTime object
func MakeFloatTime(val float64, time time.Time) FloatTime {
	return FloatTime{
		val:  val,
		time: time}
}

// Time get the time of a timed value
func (x *FloatTime) Time() time.Time {
	return x.time
}

// Val get the value of a timed value
func (x *FloatTime) Val() float64 {
	return x.val
}

// String implements the Stringer interface
func (x *FloatTime) String() string {
	return fmt.Sprintf("{val:%f, time:%s}", x.val, x.time.Format(TIME_FORMAT))
}
