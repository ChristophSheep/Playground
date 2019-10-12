package brain

import "fmt"

type Connect interface {
	AddInput(ch chan int, weight int)
	AddOutput(ch chan int)
	Update()
	Name() string
}

func ConnectBy(from, to Connect, weight int) {
	ch := make(chan int)

	from.AddOutput(ch)
	to.AddInput(ch, weight)

	from.Update()
	to.Update()

	fmt.Println(fmt.Sprintf("cell '%s' connected with '%s'", from.Name(), to.Name()))

}
