package prompts

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"whitebox/testing/abtest"
	"whitebox/testing/abtest/detect"

	"github.com/google/uuid"
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
	fmt.Println()
	fmt.Println("====================================")
	fmt.Println("            RESULT")
	fmt.Println("====================================")

	winner, ok := abtest.PickWinner(results)
	if !ok {
		fmt.Println("No results")
		return
	}

	score, _ := abtest.Score(winner.Metrics)

	fmt.Printf("🏆 Winner: %s\n", winner.Name)
	fmt.Printf("   Score: %d [%s]\n", score, abtest.ScoreTitle(score))

	sorted := abtest.Compare(results)

	fmt.Println()
	fmt.Println("RANKING:")

	for i, r := range sorted {
		score, _ := abtest.Score(r.Metrics)

		fmt.Printf("%d. %s → %d\n",
			i+1,
			r.Name,
			score,
		)
	}

	fmt.Println()
	fmt.Println("====================================")
	fmt.Println("            SAVE")
	fmt.Println("====================================")

	id := uuid.New().String()
	dir := "./ab_testing_results"

	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Println("save error:", err)
		return
	}

	type saved struct {
		Name    string         `json:"name"`
		Output  string         `json:"output"`
		Error   string         `json:"error,omitempty"`
		Score   int            `json:"score"`
		Title   string         `json:"title"`
		Metrics abtest.Metrics `json:"metrics"`
	}

	var out []saved

	for _, r := range results {
		score, _ := abtest.Score(r.Metrics)

		s := saved{
			Name:    r.Name,
			Output:  r.Output,
			Score:   score,
			Title:   abtest.ScoreTitle(score),
			Metrics: r.Metrics,
		}

		if r.Error != nil {
			s.Error = r.Error.Error()
		}

		out = append(out, s)
	}

	report := map[string]any{
		"id":      id,
		"created": time.Now(),
		"results": out,
	}

	path := filepath.Join(dir, id+".json")

	f, err := os.Create(path)
	if err != nil {
		fmt.Println("save error:", err)
		return
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")

	if err := enc.Encode(report); err != nil {
		fmt.Println("save error:", err)
		return
	}

	fmt.Println("saved to:", path)
}

func printMultiline(s string) {
	lines := strings.Split(s, "\n")

	for _, l := range lines {
		fmt.Println("  " + l)
	}
}
