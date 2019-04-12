package main

import (
	"fmt"
	"time"

	"github.com/mysheep/cell/cm3u8"
)

type TimeSlot struct {
	start time.Time
	end   time.Time
}

func (ts TimeSlot) String() string {
	return fmt.Sprintf("{start:'%v', end:'%v'}", ts.start, ts.end)
}

type DownloadOrder struct {
	channel  string
	timeSlot TimeSlot
	folder   string
}

func (do DownloadOrder) String() string {
	return fmt.Sprintf("{channel:'%v', time:%v, folder:'%v'}", do.channel, do.timeSlot, do.folder)
}

type DownloadItem struct {
	url    cm3u8.M3U8URL
	folder string
}
