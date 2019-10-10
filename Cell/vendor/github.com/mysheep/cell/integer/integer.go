package integer

import "fmt"

func AddOne(in <-chan int, out chan<- int) {
	for {
		val := <-in
		val = val + 1
		out <- val
	}
}

func AddAsync(in1, in2 <-chan int, out chan<- int) {

	var val1 = 0
	var val2 = 0

	calc := func() int {
		return val1 + val2
	}

	for {
		select {
		case val1 = <-in1:
			res := calc()
			out <- res
		case val2 = <-in2:
			res := calc()
			out <- res
			//default: // spinning of default clause exists
			// Nothing todo
		}
	}
}

func Lambda(in <-chan int, out chan<- int, fn func(int) int) {
	for {
		val := <-in
		val = fn(val)
		out <- val
	}
}

func Lambda2(in1 <-chan int, in2 <-chan int, out chan<- int, fn func(int, int) int) {
	for {
		x := <-in1
		y := <-in2
		z := fn(x, y)
		out <- z
	}
}

func Display(in <-chan int) {
	for {
		val := <-in
		fmt.Println("display", "val:", val)
	}
}

// https://stackoverflow.com/questions/19992334/how-to-listen-to-n-channels-dynamic-select-statement

func Aggregate(ins []chan int, agg chan int, out chan int) {

	//agg := make(chan int, 10)

	// TODO: Forwarder
	for i, in := range ins {
		go func(i int, ch chan int) {
			for val := range ch {
				agg <- val
			}
		}(i, in)
	}

	for {
		select {
		case val := <-agg:
			out <- val
		}
	}
}
