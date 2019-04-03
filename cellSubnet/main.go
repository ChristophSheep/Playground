package main

import (
	"fmt"
	"time"
)

func Counter(in <-chan int, out chan<- int, init int) {
	count := init
	for {
		<-in
		count = count + 1
		out <- count
	}
}

func Emitter(out chan<- int) {
	i := 7
	for {
		out <- i
		time.Sleep(50 * time.Millisecond)
		i++
	}
}

func adder(in1, in2 <-chan int, out chan<- int) {
	for {
		val1 := <-in1
		val2 := <-in2
		out <- val1 + val2
	}
}

func add1(in <-chan int, out chan<- int) {
	for {
		val := <-in
		out <- val + 1
	}
}

func Subnet(in1, in2 <-chan int, out chan<- int) {

	//            c = (a + b) +1
	//
	//        +----------------------+
	//        |                      |
	//        |    +---+             |
	// a ---->|-1->|   |    +----+   |
	//        |    | + |-3->| +1 |-4-|--c->
	// b ---->|-2->|   |    +----+   |
	//        |    +---+             |
	//        |                      |
	//        +----------------------+

	ch1 := make(chan int, 10)
	ch2 := make(chan int, 10)
	ch3 := make(chan int, 10)
	ch4 := make(chan int, 10)

	go adder(ch1, ch2, ch3)
	go add1(ch3, ch4)

	for {
		a := <-in1
		b := <-in2

		ch1 <- a
		ch2 <- b

		res := <-ch4

		out <- res
	}

}

func DoneEmitter(int <-chan int, done chan<- bool, fn func(int) bool) {
	for {
		val := <-int
		if fn(val) {
			done <- true
		}
	}
}

func Distributer(in <-chan int, out1 chan<- int, out2 chan<- int) {
	for {
		val := <-in
		out1 <- val
		out2 <- val
	}
}

func Display(in <-chan int, fn func(int)) {
	for {
		val := <-in
		fn(val)
	}
}

func main() {

	done := make(chan bool)
	waitUntilDone := func() {
		for {
			val := <-done
			if val {
				fmt.Println("!!done!!")
				break
			}
		}
	}

	ch1 := make(chan int, 10)
	ch2 := make(chan int, 10)
	ch3 := make(chan int, 10)
	ch4 := make(chan int, 10)
	ch5 := make(chan int, 10)
	ch6 := make(chan int, 10)
	ch7 := make(chan int, 10)
	ch8 := make(chan int, 10)
	ch9 := make(chan int, 10)

	quitFn := func(x int) bool {
		if x > 10 {
			return true
		}
		return false
	}

	printFn1 := func(val int) {
		fmt.Println("count", "val", val)
	}

	printFn2 := func(val int) {
		fmt.Println("subnet", "val", val)
	}

	// An emitter emits number into the network
	// The dispatch send the numer to display
	// and also to the counter
	// The counter count the emits values and
	// send the count into the quit cell.
	// The quit cell check the value and if
	// the value > 10 it send a true into
	// done channel.
	// The waitUntil is block until the done
	// has a value. So the programm blocks
	// while of go routines works. If a value
	// is send to done channel the programm ends.

	//								  		  0
	//                                        |
	//  +---------+     +-----------+     +---v---+     +------+
	//  | emitter o-1-->|  dispatch o-2-->| count o-4-->| quit o--->[waitUntil]
	//  +---------+     +-----------+     +-------+     +------+
	//	  				      | 3
	//                        v
	//				     +---------+
	//				     | display |
	//				     +---------+

	go Emitter(ch1)
	go Distributer(ch1, ch2, ch3)
	go Counter(ch2, ch4, 0 /*start value*/)
	go Distributer(ch4, ch5, ch6)
	go Display(ch5, printFn1)
	go Distributer(ch6, ch7, ch8)
	go Subnet(ch7, ch8, ch9)
	go Display(ch9, printFn2)
	go DoneEmitter(ch4, done, quitFn)

	waitUntilDone()
}
