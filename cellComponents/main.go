package main

import (
	"fmt"
	"time"
)

func counter(in <-chan int, out chan<- int, init int) {
	count := init
	for {
		<-in
		count = count + 1
		out <- count
	}
}

func emitter(out chan<- int) {
	i := 7
	for {
		out <- i
		time.Sleep(50 * time.Millisecond)
		i++
	}
}

func quit(int <-chan int, done chan<- bool, fn func(int) bool) {
	for {
		val := <-int
		if fn(val) {
			done <- true
		}
	}
}

func dispatch(in <-chan int, out1 chan<- int, out2 chan<- int) {
	for {
		val := <-in
		out1 <- val
		out2 <- val
	}
}

func display(in <-chan int, fn func(int)) {
	for {
		val := <-in
		fn(val)
	}
}

func waitUntil(done <-chan bool) {
	<-done
	fmt.Println("!!done!!")
}

func main() {

	done := make(chan bool)

	ch1 := make(chan int, 10)
	ch2 := make(chan int, 10)
	ch3 := make(chan int, 10)
	ch4 := make(chan int, 10)

	quitFn := func(x int) bool {
		if x > 10 {
			return true
		}
		return false
	}

	printFn := func(x int) {
		fmt.Println("display", "val", x)
	}

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

	go emitter(ch1)
	go dispatch(ch1, ch2, ch3)
	go display(ch3, printFn)
	go counter(ch2, ch4, 0 /*start value*/)
	go quit(ch4, done, quitFn)

	waitUntil(done)
}
