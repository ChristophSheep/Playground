package strings

import (
	"bytes"
	"fmt"
	"strings"
)

func Example1() {
	buf := make([]byte, 3, 3)
	r := strings.NewReader("Das ist ein String.")
	var n int
	var err error
	for err == nil {
		// read 3 bytes into buffer from reader
		n, err = r.Read(buf)
		fmt.Println("n", n, "bytes read to buffer", buf)
	}
}

func Example2() {
	str := "日本\x80語"
	var r rune = '\u672C' // 本 rune is a char in other languages
	index := strings.IndexRune(str, r)
	fmt.Println("index:", index)
	for pos, r := range str {
		fmt.Printf("%#U at pos %d\n", r, pos)
	}
}

func Example3() {
	str := "日本\x80語" //UTF-8
	r := strings.NewReader(str)
	// Read bytes of an UTF-8 string
	buf := make([]byte, 30)
	n, err := r.Read(buf)

	if err == nil {
		buf = buf[0:n]
		fmt.Println(buf, "len(buf)", len(buf))
	}

	// A rune is an int32
	runes := bytes.Runes(buf)
	fmt.Println(runes)
}
