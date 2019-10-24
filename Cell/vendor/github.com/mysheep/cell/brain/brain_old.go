package brain

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mysheep/timed"
)

// ----------------------------------------------------------------------------
// Cell parts
// ----------------------------------------------------------------------------

func Synapse(weight *float64, in <-chan timed.SignalTime, out chan<- timed.FloatTime) func() {

	for {
		signal := <-in
		val := timed.MakeFloatTime(0.0, signal.Time())
		if signal.Val() {
			val = timed.MakeFloatTime(*weight, signal.Time())
		}
		out <- val
	}
}

var (
	THRESHOLD = 10.0
)

func Body(agg <-chan timed.FloatTime, out chan<- timed.SignalTime, threshold float64) {

	sum := 0.0

	one := timed.MakeSignalTime(true, time.Now())
	zero := timed.MakeSignalTime(false, time.Now())

	for {
		select {
		case val := <-agg:
			sum = sum + val.Val()
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

func Axon(in <-chan timed.SignalTime, outs []chan timed.SignalTime) {
	for {
		val := <-in
		for _, out := range outs {
			out <- val
		}
	}
}

func Display(in <-chan timed.SignalTime, text string) {
	for {
		x := <-in
		fmt.Println(getNow(), "-", text, x.String())
	}
}

var (
	filenameTemplate = "cell_data_%s.txt"
)

func Writer(in <-chan timed.SignalTime, name string) {

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
		s := fmt.Sprintf("%0.3f\t%t\n", getDelta(first), val.Val()) // TODO
		s = strings.Replace(s, ".", ",", -1)
		file.WriteString(s)
	}
}
