package messages

import (
	"fmt"
)

func OutputLimit(str string, limit int) string {
	runes := []rune(str)

	if len(runes) > limit {
		visible := string(runes[:limit])
		rest := len(runes) - limit
		return fmt.Sprintf("%s...[+%d]", visible, rest)
	}

	return str
}
func LimitArgs(args map[string]string, limit int) map[string]string {
	out := make(map[string]string, len(args))

	for k, v := range args {
		out[k] = OutputLimit(v, limit)
	}

	return out
}
