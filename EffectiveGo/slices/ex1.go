package slices

import (
	"fmt"
	"os"
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
