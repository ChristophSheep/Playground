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
	outConnected chan bool
	in           chan int
	fn           func(int) int
	out          chan int
}

func add1(x int) int {
	return x + 1
}

func setIn(c *cell, in chan int) {
	c.in = in
	c.inConnected <- true
	fmt.Printf("setIn  %p %v\n", c, c)
}
func setOut(c *cell, out chan int) {
	c.out = out
	c.outConnected <- true
	fmt.Printf("setOut %p %v\n", c, c)
}

func connect(src *cell, dest *cell) {
	ch := make(chan int, 10)
	setOut(src, ch)
	setIn(dest, ch)
}

func cellRun(c *cell) {
	<-c.inConnected  // wait until input connected
	<-c.outConnected // wait until output connected
	for {
		val := c.fn(<-c.in)
		c.out <- val
	}
}

func makeCell(name string, fn func(int) int) *cell {
	ci := make(chan bool, 1)
	co := make(chan bool, 1)
	c := cell{name: name, inConnected: ci, outConnected: co, fn: fn}
	go cellRun(&c)
	return &c
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

	c1 := makeCell("cell1", add1)
	c2 := makeCell("cell2", add1)

	//
	//    [ c1 ]     [ c2 ]
	//

	fmt.Printf("address cell1 %p\n", c1)
	fmt.Printf("address cell2 %p\n", c2)

	go cellDisplay(out)

	setOut(c2, out)
	setIn(c1, in)

	// in                  out
	// -->[ c1 ]     [ c2 ]-->
	//

	connect(c1, c2)

	// in                   out
	// -->[ c1 ] ----> [ c2 ]-->
	//

	in <- 1
	in <- 2
	in <- 3

	fmt.Println("wait 5 secs ...")
	time.Sleep(5 * time.Second)
}
