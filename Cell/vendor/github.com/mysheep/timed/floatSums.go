package timed

import (
	"fmt"
	"time"
)

type FloatSums struct {
	sums   map[time.Time]float64
	maxAge int
}

// MakeFloatSums make a sums structure with timestamp sums
// in creates a hashmap with type of key time and values of type float64
// map[time]float64. If some add a new val if look for the entry with
// the same time stamp and add to this hashmap to create the sum of all
// added value with the same timestamp. To simulate that each values
// are added quasi at the same time.
func MakeFloatSums(maxAgeSeconds int) FloatSums {
	return FloatSums{
		sums:   map[time.Time]float64{},
		maxAge: maxAgeSeconds}
}

// Add add the float64 value with timestamp
// to the sum with the same timestamp
func (c *FloatSums) Add(t FloatTime) {
	if sum, ok := c.sums[t.time]; ok {
		c.sums[t.time] = sum + t.val
	} else {
		c.sums[t.time] = t.val
		c.removeOld()
	}
}

// Adds adds more then one value
func (c *FloatSums) Adds(ts ...FloatTime) {
	for _, t := range ts {
		c.Add(t)
	}
}

// Sums return the sum of a given timestamp t
func (c *FloatSums) Sum(t time.Time) (float64, bool) {
	val, ok := c.sums[t]
	return val, ok
}

// ResetSum reset the sum of a givem timestamp t
func (c *FloatSums) ResetSum(t time.Time) {
	_, ok := c.sums[t]
	if ok {
		c.sums[t] = 0.0
	}
}

// isOld check if the timestamp is older then a maxAge and now
func (c *FloatSums) isOld(t time.Time) bool {
	delta := time.Duration(c.maxAge) * time.Second
	return time.Now().Sub(t) > delta
}

// removeOld removes all to old entries in the hashmap of timestamped sums
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
