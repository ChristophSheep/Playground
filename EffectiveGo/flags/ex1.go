package flags

import (
	"flag"
	"fmt"
)

var (
	size int
	path string
)

func init() {
	fmt.Println("init of packages flags called")
	flag.IntVar(&size, "size", 32, "size of images")
	flag.StringVar(&path, "path", "./images/", "path to images")
}

func Example1() {
	flag.Parse()

	fmt.Println("parsed flags are:")
	// Visit all flags parsed and print them
	flag.Visit(
		func(f *flag.Flag) {
			fmt.Println("", "-", "value of flag", "'"+f.Name+"'", "is", f.Value, ", the usage is", "'"+f.Usage+"'")
		})

	// Return non-flag arguments
	if len(flag.Args()) > 0 {
		fmt.Println("non-flag arguments:")
	}
	for _, flag := range flag.Args() {
		fmt.Println("", "-", flag)
	}
}
