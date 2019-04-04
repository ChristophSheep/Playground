package integer

import "fmt"

func Distributor(in <-chan int, out1 chan<- int, out2 chan<- int) {
	for {
		val := <-in
		out1 <- val
		out2 <- val
	}
}

func Add1(in <-chan int, out chan<- int) {
	for {
		val := <-in
		val = val + 1
		out <- val
	}
}

func Add2Async(in1, in2 <-chan int, out chan<- int) {

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
		default:
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
