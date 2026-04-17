package embedded_prompts

import "fmt"

func IsNeedPlanModeV1(goal string) string {
	return fmt.Sprintf(`
			message-call: ask
			Only answer yes or no.
			
			Need planning/tools for this request?
			
			%s
	`, goal)
}
