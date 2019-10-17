package brain

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// ----------------------------------------------------------------------------
// Connect interface between cells
// ----------------------------------------------------------------------------

type InputConnector interface {
	InputConnect(ch chan IntTime, weight float64)
}

type OutputConnector interface {
	OutputConnect(ch chan IntTime)
}

type Namer interface {
	Name() string
}

func ConnectBy(out OutputConnector, in InputConnector, weight float64) {

	connection := make(chan IntTime)
	out.OutputConnect(connection)
	in.InputConnect(connection, weight)

	//fmt.Println(fmt.Sprintf("cell '%s' connected to '%s'", out.(Namer).Name(), in.(Namer).Name()))
}

func getNow() string {
	return time.Now().Format("15:04:05.000")
}

// ----------------------------------------------------------------------------
// Cell parts
// ----------------------------------------------------------------------------

func Synapse(weight *float64, in <-chan IntTime, out chan<- FloatTime) func() {

	for {
		signal := <-in
		val := FloatTime{val: 0.0, time: signal.time}
		if signal.val > 0 {
			val = FloatTime{val: *weight, time: signal.time}
		}
		out <- val
	}
}

var (
	THRESHOLD = 10.0
)

func Body(agg <-chan float64, out chan<- int, threshold float64) {

	sum := 0.0

	for {
		select {
		case val := <-agg:
			sum = sum + val
			if sum > threshold {
				for ; sum > threshold; sum = sum - threshold {
					out <- 1
				}
			} else {
				out <- 0
			}
		}
	}
}

func Axon(in <-chan int, outs []chan int) {
	//fmt.Printf("Axon 2 %p \n", &outs)
	for {
		val := <-in
		for _, out := range outs {
			out <- val
		}
	}
}

func Display(in <-chan IntTime, text string) {
	for {
		x := <-in
		fmt.Println(getNow(), "-", text, "val:", x.val, "time:", x.time)
	}
}

var (
	filenameTemplate = "cell_data_%s.txt"
)

func Writer(in <-chan int, name string) {

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
		s := fmt.Sprintf("%0.3f\t%d\n", getDelta(first), val)
		s = strings.Replace(s, ".", ",", -1)
		file.WriteString(s)
	}
}
