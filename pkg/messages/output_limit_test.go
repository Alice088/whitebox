package messages

import (
	"fmt"
	"strings"
)

func OutputLimit(str string, limit int) string {
	lines := strings.Split(str, "\n")

	if len(lines) > limit {
		visible := strings.Join(lines[:limit], "\n")
		rest := len(lines) - limit
		return fmt.Sprintf("%s\n...[+%d lines]", visible, rest)
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
