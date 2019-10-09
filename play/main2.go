package main

type Generic interface {
}

type Adder interface {
	add(y string) string
	len() int
}

type Object interface {
}

type Box interface {
	getVal() Object
	setVal(val Object)
}

type INT int
type FLOAT float64

func (x INT) len() int {
	return 1
}
func (x FLOAT) len() int {
	return 2
}

func (x INT) add(y string) string {
	x = x + 1
	return y
}

func (x FLOAT) add(y string) string {
	x = x + 1
	return y
}

func f(x int) func(int) func(int) int {
	return func(y int) func(int) int {
		return func(z int) int {
			return x + y + z
		}
	}
}
