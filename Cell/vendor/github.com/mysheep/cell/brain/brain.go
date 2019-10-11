package brain

import (
	"fmt"
	"time"
)

// TODO: see integer.Aggregate to work with a list of input channels

func Synapse(weight int, in <-chan int, out chan<- int) func() {

	for {
		signal := <-in
		val := 0
		if signal > 0 {
			val = weight
		}
		out <- val
	}
}

var (
	MUS = 10
)

func Cell(ins []chan int, out chan<- int) {

	agg := make(chan int)

	for i, in := range ins {
		go func(i int, ch chan int) {
			for val := range ch {
				agg <- val
			}
		}(i, in)
	}

	sum := 0

	for {
		select {
		case val := <-agg:
			sum = sum + val
			//fmt.Println("body sum:", sum)
			for ; sum > MUS; sum = sum - MUS {
				out <- 1
			}
		}
	}
}

func Axon(in <-chan int, outs []chan int) {
	for {
		val := <-in
		for _, out := range outs {
			out <- val
		}
	}
}

func Display(in <-chan int, text string) {
	for {
		val := <-in
		fmt.Println(time.Now().Format("15:04:05.000"), text, "val:", val)
	}
}
