package main

import (
	"fmt"
	"time"
)

func producer(n int, ch chan<- int) {
	for i := 0; i < n; i++ {
		ch <- i
		time.Sleep(200.0 * time.Millisecond)
	}

}

func consumer(n int, ch <-chan int, done chan int) {
	time.Sleep(2 * time.Second)
	for i := 0; i < n; i++ {
		i := <-ch // block until data come
		fmt.Printf("consumer receive val: %4d channel len: %d \n", i, len(ch))
		time.Sleep(50 * time.Millisecond)
	}
	done <- 0
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
	return x * x
}

func (c *cell) setIn(in chan int) {
	c.in = in
	c.inConnected <- true
	fmt.Printf("setIn  %p %v\n", c, c)
}

func (c *cell) setOut(out chan int) {
	c.out = out
	c.outConnected <- true
	fmt.Printf("setOut %p %v\n", c, c)
}

func (src *cell) connectTo(dest *cell) {
	ch := make(chan int, 10)
	src.setOut(ch)
	dest.setIn(ch)
}

func (c *cell) cellRun() {
	<-c.inConnected  // wait until input connected
	<-c.outConnected // wait until output connected
	for {
		val := c.fn(<-c.in)
		c.out <- val
	}
	// TODO: Disconnect stop running??
}

func makeCell(name string, fn func(int) int) *cell {
	inConnectedSignal := make(chan bool, 1)
	outConnectedSignal := make(chan bool, 1)
	c := cell{name: name, inConnected: inConnectedSignal, outConnected: outConnectedSignal, fn: fn}
	go c.cellRun()
	return &c
}

func main() {

	const N = 7

	in := make(chan int, 10)
	out := make(chan int, 10)

	c1 := makeCell("cell1", add1)
	c2 := makeCell("cell2", add1)

	//
	//    [ c1 ]     [ c2 ]
	//

	fmt.Printf("address cell1 %p\n", c1)
	fmt.Printf("address cell2 %p\n", c2)

	c2.setOut(out)
	c1.setIn(in)

	// in                  out
	// -->[ c1 ]     [ c2 ]-->
	//

	c1.connectTo(c2)

	// in                   out
	// -->[ c1 ] ----> [ c2 ]-->
	//

	done := make(chan int)

	// produce some input
	go producer(N, in)
	// consome some output
	go consumer(N, out, done)

	// Wait until all N receiver by consumer
	<-done
}
