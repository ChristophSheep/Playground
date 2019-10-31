package main

import "fmt"

func filter(xs []int, m int, fn func(int) bool) []int {
	ys := make([]int, 0, m)
	for _, x := range xs {
		if fn(x) {
			ys = append(ys, x)
		}
	}
	return ys
}

func sort(xs []int) []int {
	if len(xs) <= 1 {
		return xs
	}

	M := len(xs) / 2
	pivot := xs[M]
	left := filter(xs, M, func(x int) bool { return x < pivot })
	right := filter(xs, M, func(x int) bool { return x > pivot })
	return append(append(sort(left), pivot), sort(right)...)
}

func main() {
	xs := []int{6, 5, 2, 4, 1, 7, 5, 4, 3, 8}

	const N = 3
	fmt.Println(xs[:N]) // first 3
	fmt.Println(xs[N:]) // all others after first 3

	fmt.Println(sort(xs))
}
