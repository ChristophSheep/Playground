package main

import "fmt"

func repeater(in <-chan int, out1 chan<- int, out2 chan<- int) {
	for {
		val := <-in
		out1 <- val
		out2 <- val
	}
}

func divide(in <-chan int, out chan<- int, fn func(int) int) {
	for {
		val := <-in
		res := fn(val)
		out <- res
	}
}

func input(feedback <-chan int, out chan<- int, val int) {
	fmt.Println("input", "send", "val =", val, "into network")
	out <- val
	for {
		val := <-feedback
		out <- val
	}
}

func display(in <-chan int) {
	for {
		val := <-in
		fmt.Println("display", "print", "val =", val)
	}
}

func check(in <-chan int, out chan<- int, cond func(int) bool, done chan<- bool) {
	for {
		val := <-in
		if cond(val) {
			fmt.Println("check", "val", val, " if val < 2, send now done!!")
			done <- true
		}
		out <- val
	}
}

func switcher(in <-chan int, out1 chan<- int, out2 chan<- int, cond func(int) bool) {
	for {
		val := <-in
		if cond(val) {
			out1 <- val
		} else {
			out2 <- val
		}
	}
}

func lt(x int, y int) bool {
	if x < y {
		return true
	} else {
		return false
	}
}

func lt2(x int) bool {
	return lt(x, 2)
}

func div(x int, y int) int {
	return x / y
}

func div2(x int) int {
	return div(x, 2)
}

func main() {

	//
	//          feedback loop
	//             +----+
	//     +-----4-| <2 |-3-----+
	//     |       +----+       |
	//     4                    3
	//   +-v-+     +----+     +-^-+     [ ------- ]
	//   | 8 |-1-->| /2 |-2---| R |-4-->[ display ]
	//   +---+     +----+     +---+     [ ------- ]
	//                                      | |
	//                                    /-----/
	//
	//   input      network   repeater  output

	ch1 := make(chan int, 10)
	ch2 := make(chan int, 10)
	ch3 := make(chan int, 10)
	ch4 := make(chan int, 10)

	done := make(chan bool)

	go input(ch4, ch1, 8)
	go divide(ch1, ch2, div2)
	go repeater(ch2, ch3, ch4)
	go check(ch3, ch4, lt2, done)
	go display(ch4)

	<-done // wait until done
	fmt.Println("received done!! quit")
}
