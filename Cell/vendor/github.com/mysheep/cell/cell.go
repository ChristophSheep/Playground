package cell

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var (
	consolePS1 = "go-cell> "
)

func Console2(cmdFns map[string]func()) {

	fmt.Println("Press ESC button or Ctrl-C to exit this program")
	fmt.Println("Press any key to see their ASCII code follow by Enter")

	for {
		// only read single characters, the rest will be ignored!!
		consoleReader := bufio.NewReaderSize(os.Stdin, 1)
		fmt.Print(">")
		input, _ := consoleReader.ReadByte()

		ascii := input

		// ESC = 27 and Ctrl-C = 3
		if ascii == 27 || ascii == 3 || ascii == 9 {
			fmt.Println("Exiting...")
			os.Exit(0)
		}

		fmt.Println("ASCII : ", ascii)
	}

}

func Console(cmdFns map[string]func()) {

	reader := bufio.NewReader(os.Stdin)

	getCmd := func() (string, error) {
		cmd, err := reader.ReadString('\n')
		cmd = strings.Replace(cmd, "\n", "", -1)
		return cmd, err
	}

	printHelp := func() {
		for cmdFn := range cmdFns {
			fmt.Println("-", cmdFn)
		}
	}

	invokeCmd := func(cmd string, cmdFns map[string]func()) {
		_, exists := cmdFns[cmd]
		if exists {
			fnCmd := cmdFns[cmd]
			fnCmd()
		} else {
			fmt.Printf("%s: Command not found!", cmd)
			fmt.Println()
			fmt.Println()
			fmt.Println("Available commands:")
			printHelp()
		}
	}

	printPrompt := func() {
		fmt.Print(consolePS1)
	}

	for {
		printPrompt()
		cmd, _ := getCmd()
		invokeCmd(cmd, cmdFns)
	}
}
