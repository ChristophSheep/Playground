package brain

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
