package messages

import "fmt"

func PrintError(err error) {
	fmt.Printf("\033[91mError: %v\033[0m\n", err)
}

func PrintAssistant(str string) {
	fmt.Printf("\x1b[47m  \x1b[0m whitebox >  %s\n", str)
}
