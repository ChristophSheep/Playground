package ctime

import (
	"fmt"
	"time"
)

func Timer(times <-chan time.Time, signal chan<- bool) {
	for {
		t := <-times
		now := timeIn("Vienna") // is UTC
		durationToStart := t.Sub(now)
		fmt.Println("Timer", "durationToStart:", durationToStart, "now:", now, "t:", t)
		if durationToStart > 0 {
			time.Sleep(durationToStart)
			signal <- true
		}
	}
}

// see https://stackoverflow.com/questions/25318154/convert-utc-to-local-time-go

var countryTz = map[string]string{
	"Vienna": "Europe/Vienna",
	// ...
}

func timeIn(name string) time.Time {
	loc, err := time.LoadLocation(countryTz[name])
	if err != nil {
		panic(err)
	}
	return time.Now().In(loc)
}
