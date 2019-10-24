package main

import (
	"time"

	"github.com/mysheep/sums"
)

func main() {

	now := time.Now()

	v1 := sums.MakeIntTime(1, now)
	v2 := sums.MakeIntTime(2, now)
	v3 := sums.MakeIntTime(3, now)

	v4 := sums.MakeIntTime(4, now.Add(1*time.Second))
	v5 := sums.MakeIntTime(4, now.Add(1*time.Second))
	v6 := sums.MakeIntTime(4, now.Add(1*time.Second))

	sums := sums.MakeIntSums(1)
	sums.AddVals(v1, v2, v3)
	time.Sleep(500 * time.Millisecond)
	sums.AddVals(v4, v5, v6)
	sums.ShowMap()
}
