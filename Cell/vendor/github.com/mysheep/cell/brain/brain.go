package brain

import (
	"fmt"
	"os"
	"time"
)

// TODO: see integer.Aggregate to work with a list of input channels

func Synapse(weight int, in <-chan int, out chan<- int) func() {

	for {
		signal := <-in
		val := 0
		if signal > 0 {
			val = weight
		}
		out <- val
	}
}

var (
	MUS = 10
)

func Cell(ins []chan int, out chan<- int) {

	agg := make(chan int)

	for i, in := range ins {
		go func(i int, ch chan int) {
			for val := range ch {
				agg <- val
			}
		}(i, in)
	}

	sum := 0

	for {
		select {
		case val := <-agg:
			sum = sum + val

			if sum > MUS {
				for ; sum > MUS; sum = sum - MUS {
					out <- 1
				}
			} else {
				out <- 0
			}
		}
	}
}

func Axon(in <-chan int, outs []chan int) {
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
		fmt.Println(time.Now().Format("15:04:05.000"), text, "val:", val)
	}
}

var (
	filename = "cell_data.txt"
)

func Writer(in <-chan int) {

	file, err := os.Create(filename)
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
		file.WriteString(s)
	}
}
