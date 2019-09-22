package main

import (
	"testing"
	"time"
)

func TestGrapper(t *testing.T) {

	orders := make(chan DownloadOrder)

	xs := []struct {
		given DownloadOrder
		want  []string
	}{
		{
			DownloadOrder{
				channel: "orf1test",
				timeSlot: TimeSlot{
					start: time.Now(),
					end:   time.Now().Add(15 * time.Second),
				},
				folder: "test",
			},
			[]string{"1.ts", "2.ts", "3.ts", "4.ts"},
		},
	}

	go Grapper(orders)

	for _, x := range xs {
		orders <- x.given
	}
}
