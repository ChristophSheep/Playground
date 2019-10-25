package timed

import (
	"testing"
	"time"
)

type Test struct {
	vals   []FloatTime
	expSum float64
}

func TestAddVals(t *testing.T) {

	// needed
	const MaxAgeSeconds = 1
	now := time.Now()

	//define tests
	var tests = []Test{
		{
			vals: []FloatTime{
				MakeFloatTime(1.0, now),
				MakeFloatTime(1.0, now),
			},
			expSum: 2.0},
		{
			vals: []FloatTime{
				MakeFloatTime(1.0, now),
				MakeFloatTime(1.0, now.Add(1*time.Millisecond)),
				// Out of EPSILON of 2 ms
				MakeFloatTime(1.0, now.Add(3*time.Millisecond)),
				MakeFloatTime(1.0, now.Add(4*time.Millisecond)),
			},
			expSum: 2.0},
	}

	// test now
	for _, test := range tests {

		sums := MakeFloatSums(MaxAgeSeconds)

		for _, val := range test.vals {
			sums.Add(val)
		}

		if sum, ok := sums.Sum(now); ok == false || sum != test.expSum {
			t.Errorf("Sum(%v) = %f, expected: %f", now, sum, test.expSum)
		}

	}
}
