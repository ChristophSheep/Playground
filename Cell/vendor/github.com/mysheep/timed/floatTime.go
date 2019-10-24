package timed

import (
	"fmt"
	"time"
)

type FloatTime struct {
	val  float64
	time time.Time
}

func (x *FloatTime) Time() time.Time {
	return x.time
}

func (x *FloatTime) Val() float64 {
	return x.val
}

func (x *FloatTime) String() string {
	return fmt.Sprintf("{val:%f, time:%s}", x.val, x.time.Format(TIME_FORMAT))
}

func MakeFloatTime(val float64, time time.Time) FloatTime {
	return FloatTime{
		val:  val,
		time: time}
}
