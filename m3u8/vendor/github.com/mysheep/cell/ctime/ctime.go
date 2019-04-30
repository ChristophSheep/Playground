package ctime

import (
	"time"
)

func Timer(timesUTC <-chan time.Time, signal chan<- bool) {
	for {
		t := <-timesUTC

		durationToStart := t.Sub(time.Now())

		if durationToStart < 0 {
			durationToStart = 0
		}

		if durationToStart >= 0 {
			time.Sleep(durationToStart)
			signal <- true
		}
	}
}
