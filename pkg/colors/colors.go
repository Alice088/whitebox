package colors

import "fmt"

func Dim(str string) string {
	return fmt.Sprintf("\033[2m%s\033[0m", str)
}
