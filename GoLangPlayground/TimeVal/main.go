package main

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"
)

type timedVal struct {
	val  int
	t    time.Time
	used bool
}

const DELTA_TIME = 10

const EPSILON = 10

func sendOut(sum int, val int, out chan<- int) int {
	sum = sum + val
	if sum > EPSILON {
		out <- 1
		sum = sum - EPSILON
	} else {
		out <- 0
	}
	return sum
}

// TODO: Sum values and if sum > Epsilon send 1 and remove Epsolin from sum

// google: golang select over array of channels
// stackoverflow answer

var weights = map[int]int{
	1: 3,
	2: 5,
	3: 2,
}

// TODO: Dynamically add new channel at runtime

func cellAgg(ins []chan int, agg chan int) {
	for i, in := range ins {
		go func(i int, ch chan int) {
			for val := range ch { // TODO range ???
				agg <- (val + weights[i])
			}
		}(i, in)
	}

}

func cellSum(agg chan int, out chan<- int) {
	//
	// Now listen to aggregate channel
	//
	sum := 0
	for {
		select {
		case valWeighted := <-agg:
			sum = sendOut(sum, valWeighted, out)
		}
	}
}

func addZeros(s string, n int) string {
	for i := 0; i < n; i++ {
		s = s + "0"
	}
	return s
}

func getNow() string {
	s := time.Now().Format(time.RFC3339Nano)
	n := 32 - len(s)
	s = strings.Split(s, "+")[0]
	s = addZeros(s, n)
	return s
}

func display(in <-chan int) {
	for {
		val := <-in
		fmt.Println(getNow(), "- val:", val)
	}
}
func emitter(outs *[]chan int) {
	r := rand.New(rand.NewSource(99))

	var getZeroOrOne = func(r *rand.Rand) int {
		i := r.Int31n(255)
		j := int(math.Mod(float64(i), float64(2)))
		return j
	}

	var getSleepTime = func() time.Duration {
		delta := 1 + r.Int31n(5) // 0..4 + 5 = 0..9
		return time.Duration(delta) * time.Millisecond
	}

	var getRandomIndex = func(n int) int {
		idx := r.Int31n(int32(n))
		return int(idx)
	}

	for {
		idx := getRandomIndex(len((*outs)))
		(*outs)[idx] <- getZeroOrOne(r)
		time.Sleep(getSleepTime())

		fmt.Print("len outs:", len((*outs)), " ")

	}
}

func main_OLD() {

	//         +---------+
	//   A --->|    +--  |
	//   B --->|    |    |---> O
	//   C --->|  --+    |
	//         +---------+

	chA := make(chan int, 10)
	chB := make(chan int, 10)
	chC := make(chan int, 10)

	chIns := []chan int{chA, chB, chC}
	chOut := make(chan int, 10)

	agg := make(chan int, 100)
	cellAgg(chIns, agg)
	go cellSum(agg, chOut)

	go emitter(&chIns)
	go display(chOut)

	time.Sleep(500 * time.Millisecond)

	// TODO: Dynamically add new channel at runtime
	for i := 0; i < 3; i++ {
		chIns = append(chIns, make(chan int, 10))
	}
	// TODO: Update aggregate
	cellAgg(chIns, agg)
	// TODO: Old Channels still exists

	// TODO: Update cell with new input channels

	wait := 2
	fmt.Println("wait ", wait, " secs ...")
	time.Sleep(time.Duration(wait) * time.Second)
}
