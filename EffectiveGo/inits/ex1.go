package inits

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

var (
	home   = os.Getenv("HOME")
	user   = os.Getenv("USER")
	gopath = os.Getenv("GOPATH")
)

func init() {
	if user == "" {
		log.Fatal("$User user set")
	}
	if home == "" {
		home = "home/" + user
	}
	if gopath == "" {
		gopath = home + "/go"
	}

	flag.StringVar(&gopath, "gopath", gopath, "override default GOPATH")

}

func Example1() {

	flag.Parse() // parse 'gopath'-flag

	format := "%10s = '%s'\r\n"
	fmt.Printf(format, "home", home)
	fmt.Printf(format, "user", user)
	fmt.Printf(format, "gopath", gopath)

	i := 0
	j := 0
	ps := "\\|/-"
	qs := "-/|\\"
	for ; i < 60; i++ {
		fmt.Printf("\t%c %2d sec %c\r", qs[j], i, ps[j]) // progress bar
		time.Sleep(500 * time.Millisecond)
		if j++; j >= len(ps) {
			j = 0
		}
	}
	fmt.Printf("\tit took %d sec        \n", i)
}
