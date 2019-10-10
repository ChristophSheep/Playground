package brain

// TODO: see integer.Aggregate to work with a list of input channels

func makeSynapse(weight int, in <-chan bool, out chan<- int) func() {

	var w = weight

	return func() {
		for {
			signal := <-in
			val := 0
			if signal {
				val = w
			}
			out <- val
		}
	}

}

func makeCellBody(ins []chan int, out chan<- bool) {

}

func axon(in <-chan bool, outs []chan bool) {
	for {
		val := <-in
		for _, out := range outs {
			out <- val
		}
	}
}
