package main

import (
	"fmt"
)

func main() {

	size := 2
	zs := []int{0, 1, 2}
	zs = append(zs, 3, 4, 5)

	//xs := zs[:size] // car [0 1]
	//ys := zs[size:] // cdr [2 3 4 5]
	//      ( ...
	//  ..  ]
	//  0 1 2 3 4 5
	// [a b c d e f]

	for i := 0; i < len(zs); i = i + size {
		xs := zs[i : i+size]
		go func(ys []int) {
			fmt.Println(ys)
		}(xs)
	}

	//time.Sleep(1 * time.Second)

	ws := make([]int, 10)
	for i := range ws {
		ws[i] = i + 1
	}
	fmt.Println(ws)

	xs := ws
	ws[2] = 0

	fmt.Println(xs[0:5])
}
