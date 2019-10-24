package slices

import (
	"fmt"
	"os"
	"time"
)

var (
	zs = []int{1, 2, 3}
)

var (
	home   = os.Getenv("HOME")
	user   = os.Getenv("USER")
	gopath = os.Getenv("GOPATH")
)

func init() {

	fmt.Println("init() called")
	fmt.Println()
	fmt.Println("home:", home)
	fmt.Println("user:", user)
	fmt.Println("gopath:", gopath)
	fmt.Println()
}

func Example0() {
	fmt.Println("Example0:")

	xs := make([]int, 0 /*len*/, 8 /*cap*/)
	xs = append(xs, zs...)
	ys := xs[1:]

	fmt.Println("xs ..", cap(xs), xs)
	fmt.Println("ys ..", cap(ys), " ", ys)
	fmt.Println()

	xs[1] = 0
	fmt.Println("xs ..", cap(xs), xs)
	fmt.Println("ys ..", cap(ys), " ", ys)
	fmt.Println()

	// Underlying array of xs and ys are the same
	ys[1] = 2
	xs = append(xs, xs...)
	fmt.Println("xs ..", cap(xs), xs)
	fmt.Println("ys ..", cap(ys), " ", ys)
	fmt.Println()

	// if we append something to ys
	// the underlying array must change because
	// we would overwrite data in xs
	ys = append(ys, 0, 1, 2)
	fmt.Println("ys ..", cap(ys), " ", ys)
	fmt.Println()

	// We change some in xs and you
	// see ys is not the same anymore
	xs[1] = 9
	xs[2] = 3
	fmt.Println("xs ..", cap(xs), xs)
	fmt.Println("ys ..", cap(ys), " ", ys)
	fmt.Println()

	ys[0] = 0
	fmt.Println("xs ..", cap(xs), xs)
	fmt.Println("ys ..", cap(ys), " ", ys)
}

func Example1() {
	fmt.Println("Example1:")

	xs := make([]int, 0)
	xs = append(xs, 1, 2, 3)
	fmt.Println("xs ..", xs)
	fmt.Println("ys = xs[1:]")
	ys := xs[1:]
	fmt.Println("ys ..", ys)
	fmt.Println("xs = append(xs, 4, 5, 6)")
	xs = append(xs, 4, 5, 6)
	fmt.Println("xs ..", xs)
	fmt.Println("ys ..", ys)
	ys = append(ys, 7, 8, 9)
	fmt.Println("ys ..", ys)
	fmt.Println("xs[1] = 0")
	xs[1] = 0
	fmt.Println("xs ..", xs)
	fmt.Println("ys ..", ys)
}

func cdr(ys []int) []int {
	if len(ys) == 0 {
		return ys
	}
	return ys[1:]
}

func car(ys []int) []int {
	if len(ys) == 0 {
		panic("car of empty list not defined")
	}
	return ys[0:1]
}

// Keyword 'map' already reserved for hashmap
//
func maps(ys []int, fn func(int) int) []int {
	if len(ys) == 0 {
		return ys // or make([]int, 0) // or nil
	}
	xs := car(ys)
	xs[0] = fn(xs[0])
	return append(xs, maps(cdr(ys), fn)...)
}

func Example2() {
	// create a map function

	xs := make([]int, 0, 3) // 0 elements but cap 3
	xs = append(xs, 1, 2, 3)

	xs = maps(xs, func(x int) int { return x + 1 })
	fmt.Println("xs:", xs)

	//ys := [3]int{4, 5, 6} // array
	//ys = append(ys, 7) // does not work

	//ys := [...]int{4, 5, 6} // array
	//ys = append(ys, 7) // does not work

	ys := []int{4, 5, 6}
	ys = append(ys, 7) // work .. append only work with slices
	fmt.Println("ys:", ys)
}

func Example3() {
	const N = 4096

	fmt.Println("test - append", N, "channels")
	var test = func(fn func()) {
		start := time.Now()
		fn()
		duration := time.Now().Sub(start).Seconds() * 1000.0
		fmt.Printf("it took %.3f ms\n", duration)
	}

	var f0 = func() {
		ss := make([]chan int, 0, N)
		ch := make(chan int)
		for i := 0; i < N; i++ {
			ss = append(ss, ch)
		}
	}

	var f1 = func() {
		ss := make([]chan int, 0)
		for i := 0; i < N; i++ {
			ss = append(ss, make(chan int))
		}
	}

	var f2 = func() {
		ss := make([]chan int, N, N)
		for i := 0; i < N; i++ {
			ss[i] = make(chan int)
		}
	}

	var f3 = func() {
		ss := make([]chan int, N, N)
		ch := make(chan int)
		for i := 0; i < N; i++ {
			ss[i] = ch
		}
	}

	test(f0)
	test(f1)
	test(f2)
	test(f3)

	//it took 0.000113012 s
	//it took 0.00225483 s
	//it took 0.001857907 s
	//it took 0.000133485 s
}
