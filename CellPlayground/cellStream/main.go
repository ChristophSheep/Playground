package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type Command struct {
	name string
}

// emitter emits commands in to commands stream
func Emitter(cmdIn <-chan Command, cmdOut chan<- Command, emit <-chan bool) {

	doEmit := false

	go func() {
		for {
			<-emit
			doEmit = true
			time.Sleep(100 * time.Millisecond)
			if doEmit == false {
				doEmit = false
			}
		}
	}()

	for {
		cmd := <-cmdIn

		// TODO:
		fmt.Println("emitter receive command", cmd.name)
		// TODO:

		if doEmit {
			cmdEmit := Command{"emit"}
			cmdOut <- cmdEmit
		}

		cmdOut <- cmd

	}
}

func Cell(cmdIn <-chan Command, cmdOut chan<- Command) {
	for {
		cmd := <-cmdIn

		// TODO:
		fmt.Println("cell receive command", cmd.name)
		// TODO:

		time.Sleep(10 * time.Millisecond)

		cmdOut <- cmd
	}
}

func Console(done chan<- bool, emit chan<- bool) {

	const NEWLINE = '\n'
	scanner := bufio.NewScanner(os.Stdin)

	for {

		fmt.Print("command > ")
		//text, _ := reader.ReadString(NEWLINE)
		scanner.Scan()
		text := scanner.Text()
		text = strings.ToLower(text)

		switch text {
		case "ok":
			fmt.Println("Send", text, "!")
			done <- true

		case "emit":
			fmt.Println("Send", text, "!")
			emit <- true
		default:
			fmt.Println("Unknown command!")
		}

	}
}

func main() {

	done := make(chan bool)
	emit := make(chan bool, 10)

	ch0 := make(chan Command, 10)
	ch1 := make(chan Command, 10)
	ch2 := make(chan Command, 10)
	ch3 := make(chan Command, 10)

	go Emitter(ch0, ch1, emit)
	go Cell(ch1, ch2)
	go Cell(ch2, ch3)
	go Cell(ch3, ch0)

	startCmd1 := Command{"start"}
	ch0 <- startCmd1
	startCmd2 := Command{"foo"}
	ch0 <- startCmd2

	go Console(done, emit)

	<-done
}
