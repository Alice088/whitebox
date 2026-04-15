package prompts

import (
	"fmt"
	"strings"
)

func PrintReport(results []Result) {
	fmt.Println("====================================")
	fmt.Println("         PROMPT TEST REPORT")
	fmt.Println("====================================")

	for i, r := range results {
		fmt.Println()
		fmt.Println("────────────────────────────────────")
		fmt.Printf("CASE %d: %s\n", i+1, r.Name)
		fmt.Println("────────────────────────────────────")

		if r.Error != nil {
			fmt.Println("STATUS: ERROR")
			fmt.Println("DETAILS:")
			fmt.Println(r.Error)
			continue
		}

		fmt.Println("STATUS: OK")
		fmt.Println()

		fmt.Println("METRICS:")
		fmt.Println("--------")
		fmt.Printf("  Events:     %d\n", r.Metrics.EventCalls)
		fmt.Printf("  Steps:      %d\n", r.Metrics.Steps)
		fmt.Printf("  ToolCalls:  %d\n", r.Metrics.ToolCalls)
		fmt.Printf("  Errors:     %d\n", r.Metrics.Errors)
		fmt.Printf("  Duration:   %s\n", r.Metrics.Duration)
		fmt.Println("--------")
		fmt.Println("")

		fmt.Println("OUTPUT:")
		fmt.Println("--------")

		printMultiline(r.Output)

		fmt.Println("--------")
	}

	fmt.Println()
	fmt.Println("====================================")
}

func printMultiline(s string) {
	lines := strings.Split(s, "\n")

	for _, l := range lines {
		fmt.Println("  " + l)
	}
}
