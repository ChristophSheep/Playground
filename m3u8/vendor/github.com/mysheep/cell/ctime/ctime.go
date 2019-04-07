package ctime

import (
	"time"
)

func Timer(times <-chan time.Time, signal chan<- bool) {
	for {
		t := <-times
		durationToStart := t.Sub(time.Now())

		if durationToStart > 0 {
			time.Sleep(durationToStart)
			signal <- true
		}
	}
}
