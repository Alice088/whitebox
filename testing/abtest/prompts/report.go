package prompts

import (
	"fmt"
	"strings"
	"whitebox/testing/abtest"
	"whitebox/testing/abtest/detect"
)

func PrintReport(results []abtest.Result) {
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
		fmt.Printf("  Events:            %d\n", r.Metrics.EventCalls)
		fmt.Printf("  Steps:             %d\n", r.Metrics.Steps)
		fmt.Printf("  ToolCallsHistory:  %s\n", abtest.StringToolCallsHistory(r.Metrics.ToolsCallsHistory))
		fmt.Printf("  ToolCalls:  		%s\n", abtest.StringToolCalls(r.Metrics.ToolsCalls))
		fmt.Printf("  Errors:     		%d\n", r.Metrics.Errors)
		fmt.Printf("  Duration:   		%s\n", r.Metrics.Duration)

		toolRepeatCount := detect.ToolRepeat(r.Metrics.ToolsCallsHistory)
		isToolRepeat, toolRepeatTitle := func() (bool, string) {
			if toolRepeatCount > 0 {
				return true, "--> [warning!]"
			}
			return false, ""
		}()
		fmt.Printf("  ToolRepeat: 		%t %s\n", isToolRepeat, toolRepeatTitle)

		toolLoop := detect.ToolLoop(r.Metrics.ToolsCallsHistory)
		fmt.Printf("  ToolLoop:   		%t %s\n", toolLoop, func() string {
			if toolLoop {
				return "--> [maybe ok]"
			}
			return ""
		}())

		fmt.Println("--------")
		score, breakpoints := abtest.Score(r.Metrics)
		fmt.Printf("└─> Score: %dpts [%s]\n", score, abtest.ScoreTitle(score))
		fmt.Println("  └─> Details:")

		for k, v := range breakpoints {
			fmt.Printf("    - %s: %d\n", k, v)
		}
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
