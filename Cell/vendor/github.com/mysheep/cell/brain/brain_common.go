package brain

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// ----------------------------------------------------------------------------
// Cell parts
// ----------------------------------------------------------------------------

func Synapse(weight *float64, in <-chan SignalTime, out chan<- FloatTime) func() {

	for {
		signal := <-in
		val := FloatTime{val: 0.0, time: signal.time}
		if signal.val {
			val = FloatTime{val: *weight, time: signal.time}
		}
		out <- val
	}
}

var (
	THRESHOLD = 10.0
)

func Body(agg <-chan FloatTime, out chan<- SignalTime, threshold float64) {

	sum := 0.0

	one := SignalTime{val: true, time: time.Now()}
	zero := SignalTime{val: false, time: time.Now()}

	for {
		select {
		case val := <-agg:
			sum = sum + val.val
			if sum > threshold {
				for ; sum > threshold; sum = sum - threshold {
					out <- one
				}
			} else {
				out <- zero // TODO: ??
			}
		}
	}
}

func Axon(in <-chan SignalTime, outs []chan SignalTime) {
	for {
		val := <-in
		for _, out := range outs {
			out <- val
		}
	}
}

func Display(in <-chan SignalTime, text string) {
	for {
		x := <-in
		fmt.Println(getNow(), "-", text, x.String())
	}
}

var (
	filenameTemplate = "cell_data_%s.txt"
)

func Writer(in <-chan SignalTime, name string) {

	file, err := os.Create(fmt.Sprintf(filenameTemplate, name))
	if err != nil {
		return
	}

	flag := true
	var first time.Time

	var setFirst = func(flag *bool, first *time.Time) {
		if *flag {
			*first = time.Now()
			*flag = false
		}
	}

	var getDelta = func(first time.Time) float64 {
		return time.Now().Sub(first).Seconds()
	}

	for val := range in {
		setFirst(&flag, &first)
		s := fmt.Sprintf("%0.3f\t%t\n", getDelta(first), val.val) // TODO
		s = strings.Replace(s, ".", ",", -1)
		file.WriteString(s)
	}
}
