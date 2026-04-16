package messages

import (
	"fmt"
	"strings"
)

func StringArgs(args map[string]string) string {
	if len(args) == 0 {
		return ""
	}

	var parts []string

	for k, v := range args {
		parts = append(parts, fmt.Sprintf("%s:%s", k, OutputLimit(v, 30)))
	}

	return strings.Join(parts, ", ")
}
