package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func Adder(in <-chan int, out chan<- int) {
	for {
		val := <-in
		val = val + 1
		out <- val
	}
}

func Display(in <-chan int) {
	for {
		val := <-in
		fmt.Println("display:", val)
		printCmd()
	}
}

func printCmd() {
	fmt.Print("command > ")
}

// TODO: MOVE TO LIBRARY
func Console(done chan<- bool, cmds map[string]func()) {

	printAvailableCmds := func() {
		fmt.Println("Available commands:")
		for cmd := range cmds {
			fmt.Println("-", cmd)
		}
	}

	printUnknownCmd := func() {
		fmt.Println("Unknown command! Type 'help' for available commands!")
	}

	findAndInvokeCommand := func(text string) bool {
		_, cmdExists := cmds[text]
		if cmdExists {
			cmds[text]()
			return true
		} else {
			return false
		}
	}

	const NEWLINE = '\n'
	scanner := bufio.NewScanner(os.Stdin)

	for {

		printCmd()

		scanner.Scan()
		text := scanner.Text()
		cmd := strings.ToLower(text)

		if cmd == "help" {
			printAvailableCmds()
		} else {
			if findAndInvokeCommand(cmd) == false {
				printUnknownCmd()
			}
		}

	}
}

func emit(in chan<- int, out <-chan int) {
	start := time.Now()

	in <- 1
	val := <-out

	end := time.Now()
	elapsed := end.Sub(start)
	fmt.Println("elapsed:", elapsed, "val:", val)
}

func inputN() int {
	fmt.Print("N ? ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	text := strings.ToLower(scanner.Text())

	i, err := strconv.Atoi(text)
	if err == nil {
		fmt.Println("integer value is:", i)
	} else {
		fmt.Println("value is not a int")
	}
	return i
}

func main() {

	done := make(chan bool)
	waitUntilDone := func() {
		<-done
	}

	N := inputN()
	//const N = 10000 // e.g. 1 Mio take 1.6sec, 10.000 takes 13ms

	// http://localhost:6060/doc/effective_go.html#arrays
	chs := make([]chan int, N)
	for i := 0; i < N; i++ {
		chs[i] = make(chan int)
	}

	//  0     1             N-1
	// -->[ ]--> .... -->[ ]-->
	//

	last := N - 1
	first := 0

	for i := first; i < last; i++ {
		go Adder(chs[i], chs[i+1])
	}

	cmds := map[string]func(){
		"emit": func() { emit(chs[first], chs[last]) },
		"ok":   func() { done <- true },
	}
	go Console(done, cmds)

	waitUntilDone()
}
