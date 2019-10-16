package sums

import (
	"fmt"
	"time"
)

type IntTime struct {
	val  int
	time time.Time
}

type IntSums struct {
	sums   map[time.Time]int
	maxAge int
}

func MakeIntTime(val int, time time.Time) IntTime {
	return IntTime{
		val:  val,
		time: time}
}

func (c *IntSums) AddVal(t IntTime) {
	if sum, ok := c.sums[t.time]; ok {
		c.sums[t.time] = sum + t.val
	} else {
		c.sums[t.time] = t.val
		c.removeOld()
	}
}

func (c *IntSums) AddVals(ts ...IntTime) {
	for _, t := range ts {
		c.AddVal(t)
	}
}

func (c *IntSums) isOld(t time.Time) bool {
	delta := time.Duration(c.maxAge) * time.Second
	return time.Now().Sub(t) > delta
}

func (c *IntSums) removeOld() {
	for key, _ := range c.sums {
		if c.isOld(key) {
			delete(c.sums, key)
		}
	}
}

func (c *IntSums) ShowMap() {
	fmt.Printf("map: %v \n", c.sums)
}

func MakeIntSums(maxAgeSeconds int) *IntSums {
	return &IntSums{
		sums:   map[time.Time]int{},
		maxAge: maxAgeSeconds}
}
