package cstring

import "time"

func Switch(on bool, in <-chan string, out chan<- string) {

	for {
		val := <-in
		if on {
			out <- val
		}
	}

}

func Repeater(interval int, in <-chan string, out chan<- string) {

	val := ""
	valIsSet := false

	go func() {
		for {
			val = <-in
			valIsSet = true
		}
	}()

	for {
		time.Sleep(time.Duration(interval) * time.Second)
		if valIsSet {
			out <- val
		}
	}
}
