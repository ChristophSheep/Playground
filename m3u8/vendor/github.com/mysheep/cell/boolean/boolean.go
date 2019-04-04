package boolean

func Distributor(in <-chan bool, out1 chan<- bool, out2 chan<- bool) {
	for {
		val := <-in
		out1 <- val
		out2 <- val
	}
}
