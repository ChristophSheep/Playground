package boolean

func Distributor(in <-chan bool, out1 chan<- bool, out2 chan<- bool) {
	for {
		val := <-in
		out1 <- val
		out2 <- val
	}
}

func Forwarder(in <-chan bool, out chan<- bool) {
	for {
		val := <-in
		out <- val
	}
}
