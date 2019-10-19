package brain

import (
	"fmt"
	"time"
)

type FloatTime struct {
	val  float64
	time time.Time
}

func (x *FloatTime) String() string {
	return fmt.Sprintf("{val:%f, time:%s}", x.val, x.time.Format(TIME_FORMAT))
}

type FloatSums struct {
	sums   map[time.Time]float64
	maxAge int
}

func MakeFloatTime(val float64, time time.Time) FloatTime {
	return FloatTime{
		val:  val,
		time: time}
}

func (c *FloatSums) AddVal(t FloatTime) {
	if sum, ok := c.sums[t.time]; ok {
		c.sums[t.time] = sum + t.val
	} else {
		c.sums[t.time] = t.val
		c.removeOld()
	}
}

func (c *FloatSums) AddVals(ts ...FloatTime) {
	for _, t := range ts {
		c.AddVal(t)
	}
}

func (c *FloatSums) isOld(t time.Time) bool {
	delta := time.Duration(c.maxAge) * time.Second
	return time.Now().Sub(t) > delta
}

func (c *FloatSums) getSum(t time.Time) (float64, bool) {
	val, ok := c.sums[t]
	return val, ok
}

func (c *FloatSums) resetSum(t time.Time) {
	_, ok := c.sums[t]
	if ok {
		c.sums[t] = 0.0
	}
}

func (c *FloatSums) removeOld() {
	for key, _ := range c.sums {
		if c.isOld(key) {
			delete(c.sums, key)
		}
	}
}

func (c *FloatSums) ShowMap() {
	fmt.Printf("map: %v \n", c.sums)
}

func MakeFloatSums(maxAgeSeconds int) FloatSums {
	return FloatSums{
		sums:   map[time.Time]float64{},
		maxAge: maxAgeSeconds}
}
