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
	InputConnect(ch chan int, weight float64)
}

type OutputConnector interface {
	OutputConnect(ch chan int)
}

type Namer interface {
	Name() string
}

func ConnectBy(out OutputConnector, in InputConnector, weight float64) {

	connection := make(chan int)
	out.OutputConnect(connection)
	in.InputConnect(connection, weight)

	//fmt.Println(fmt.Sprintf("cell '%s' connected to '%s'", out.(Namer).Name(), in.(Namer).Name()))
}

// ----------------------------------------------------------------------------
// Cell parts
// ----------------------------------------------------------------------------

func Synapse(weight float64, in <-chan int, out chan<- float64) func() {

	for {
		signal := <-in
		val := float64(0.0)
		if signal > 0 {
			val = weight
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

func Display(in <-chan int, text string) {
	for {
		val := <-in
		fmt.Println(time.Now().Format("15:04:05.000"), "-", text, "val:", val)
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
