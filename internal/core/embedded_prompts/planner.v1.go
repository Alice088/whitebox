package embedded_prompts

import "fmt"

func PlannerV1(goal string) string {
	return fmt.Sprintf(`
			You are an expert planner.
			Goal: %s
			Rules:
			- 3 to 7 steps
			- short concrete steps
			- no explanations
			- no markdown
	`, goal)
}
