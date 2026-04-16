package messages

import (
	"fmt"
	"strings"
)

func OutputLimit(str string, limit int) string {
	str = Flat(str)

	if len(str) > limit {
		return str[:limit] + fmt.Sprintf("...[+%d words]", len(str[limit:]))
	}
	return str
}

func Flat(str string) string {
	return strings.Join(strings.Fields(str), " ")
}

func LimitArgs(args map[string]string, limit int) map[string]string {
	out := make(map[string]string, len(args))

	for k, v := range args {
		out[k] = OutputLimit(v, limit)
	}

	return out
}
