package integer

import "fmt"

func AddOne(in <-chan int, out chan<- int) {
	for {
		val := <-in
		val = val + 1
		out <- val
	}
}

func AddAsync(in1, in2 <-chan int, out chan<- int) {

	var val1 = 0
	var val2 = 0

	calc := func() int {
		return val1 + val2
	}

	for {
		select {
		case val1 = <-in1:
			res := calc()
			out <- res
		case val2 = <-in2:
			res := calc()
			out <- res
		}
	}
}

func Lambda(in <-chan int, out chan<- int, fn func(int) int) {
	for {
		val := <-in
		val = fn(val)
		out <- val
	}
}

func Lambda2(in1 <-chan int, in2 <-chan int, out chan<- int, fn func(int, int) int) {
	for {
		x := <-in1
		y := <-in2
		z := fn(x, y)
		out <- z
	}
}

func Display(in <-chan int) {
	for {
		val := <-in
		fmt.Println("display", "val:", val)
	}
}

// https://stackoverflow.com/questions/19992334/how-to-listen-to-n-channels-dynamic-select-statement

func Aggregate(ins []chan int, agg chan int, out chan int) {

	for i, in := range ins {
		go func(i int, ch chan int) {
			for val := range ch {
				agg <- val
			}
		}(i, in)
	}

	for {
		select {
		case val := <-agg:
			out <- val
		}
	}

}

func MakeAgg(ins *[]chan int, out chan int) (func(), func(), func()) {

	agg := make(chan int, 10)
	exit := make(chan int)
	exited := make(chan int)

	var inFn = func(i int, ch chan int) {
	Loop:
		for {
			select {
			case val := <-ch:
				agg <- val
			case <-exit:
				fmt.Println("exit break", i, "agg fn")
				break Loop
			}
		}
		exited <- i
		fmt.Println("exited", i, "agg fn")
	}

	var exitFn = func() {
		N := len((*ins))

		// Send each one signal to exit
		for i := 0; i < N; i++ {
			fmt.Println("send", i, "to exit")
			exit <- i
		}

		// Wait until all N exited
		for j := 0; j < N; j++ {
			fmt.Println("wait", j, "to exit")
			<-exited
		}

	}

	var updateAggsFn = func() {
		fmt.Println("update", len((*ins)), "ins")
		for i, in := range *ins {
			fmt.Println("create", i, "agg fn")
			go inFn(i, in)
		}
	}

	var aggFn = func() {
		for {
			select {
			case val := <-agg:
				out <- val
			}
		}
	}

	return updateAggsFn, aggFn, exitFn
}

func MakeDynAgg(ins *[]chan int, out chan int) (func(chan int), func()) {

	agg := make(chan int)

	var inFn = func(i int, ch chan int) {
		for val := range ch {
			agg <- val
		}
	}

	for i, in := range *ins {
		go inFn(i, in)
	}

	var addFn = func(in chan int) {
		i := len(*ins)
		*ins = append(*ins, in)
		go inFn(i, in)
	}

	var aggFn = func() {
		for {
			select {
			case val := <-agg:
				out <- val
			}
		}
	}

	return addFn, aggFn
}
