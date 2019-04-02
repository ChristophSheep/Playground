package main

import (
	"fmt"
	"time"
)

func producer(ch chan<- int, done chan int) {
	for i := 0; i < 7; i++ {
		ch <- i
		time.Sleep(200.0 * time.Millisecond)
	}
	done <- 0
}

func consumer(ch <-chan int) {
	time.Sleep(2 * time.Second)
	for {
		i := <-ch // block until data come
		fmt.Printf("val: %d len: %d \n", i, len(ch))
		time.Sleep(50 * time.Millisecond)
	}
}

type cell struct {
	name         string
	inConnected  chan bool
	outConnected bool
	in           chan int
	fn           func(int) int
	out          chan int
}

func add1(x int) int {
	return x + 1
}

func setIn(c cell, in chan int) {
	c.in = in
	c.inConnected <- true
}
func setOut(c cell, out chan int) {
	c.out = out
	c.outConnected = true
	fmt.Printf("setOut: %v", c)
}

func connect(src cell, dest cell) {
	ch := make(chan int, 10)
	setIn(dest, ch)
	setOut(src, ch)
}

func cellRun(c cell) {
	for {
		<-c.inConnected // wait until
		fmt.Println("connected:", c.name)
		for {
			val := c.fn(<-c.in)
			fmt.Printf("cell run: %+v address:%p", c, &c)
			fmt.Println(" - val:", val, "outConnected:", c.outConnected)
			//if c.outConnected {
			c.out <- val
			//}
		}
	}
}

func makeCell(name string, in chan int, fn func(int) int, out chan int) cell {
	connected := make(chan bool)
	c := cell{name, connected, false, in, fn, out}
	go cellRun(c)
	return c
}

//
// Create a cell and when it fully connected let it run
//

func cellDisplay(out <-chan int) {
	for {
		res := <-out
		fmt.Printf("res: %d\n", res)
	}
}

func main() {

	//done := make(chan bool)

	in := make(chan int, 10)
	out := make(chan int, 10)
	ch1 := make(chan int, 10)

	in <- 1
	in <- 2
	in <- 3

	c1 := makeCell("cell1", in, add1, ch1)
	c2 := makeCell("cell2", ch1, add1, out)

	fmt.Println()
	fmt.Printf("cell c1: %v address: %p \n", c1, &c1)
	fmt.Printf("cell c2: %v address: %p \n", c2, &c2)

	go cellDisplay(out)

	setIn(c1, in)
	setOut(c2, out)

	connect(c1, c2)

	time.Sleep(5 * time.Second)
}
