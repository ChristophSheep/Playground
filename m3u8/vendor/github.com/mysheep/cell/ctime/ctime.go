package ctime

import (
	"fmt"
	"time"
)

func Timer(timesUTC <-chan time.Time, signal chan<- bool) {
	for {
		t := <-timesUTC
		durationToStart := t.Sub(time.Now().UTC()) // MUST UTC!!! Create UTC datatype???

		fmt.Println("Timer - durationToStart:", durationToStart)

		if durationToStart < 0 {
			durationToStart = 0
		}

		if durationToStart >= 0 {
			time.Sleep(durationToStart)
			signal <- true
		}
	}
}
