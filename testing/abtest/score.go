package abtest

import "whitebox/testing/abtest/detect"

func Score(m Metrics) (int, map[string]int) {
	details := make(map[string]int)
	score := 0

	// success (если есть)
	if m.Errors == 0 {
		score += 100
		details["success"] = +100
	}

	// steps
	stepsPenalty := m.Steps * -5
	score += stepsPenalty
	details["steps"] = stepsPenalty

	// tool calls
	toolPenalty := len(m.ToolsCallsHistory) * -10
	score += toolPenalty
	details["tool_calls"] = toolPenalty

	// duration (в секундах)
	sec := m.Duration.Seconds()

	var timePenalty int

	switch {
	case sec < 5:
		timePenalty = 0
	case sec < 15:
		timePenalty = -5
	case sec < 30:
		timePenalty = -10
	case sec < 60:
		timePenalty = -20
	default:
		timePenalty = -30
	}
	score += timePenalty
	details["time"] = timePenalty

	// repeat
	repeats := detect.ToolRepeat(m.ToolsCallsHistory)
	if repeats > 0 {
		p := repeats * -10
		score += p
		details["repeat"] = p
	}

	// loop
	if detect.ToolLoop(m.ToolsCallsHistory) {
		score -= 20
		details["loop"] = -20
	}

	return score, details
}

func ScoreTitle(score int) string {
	switch {
	case score >= 90:
		return "clean as hell"
	case score >= 70:
		return "solid run"
	case score >= 50:
		return "works but meh"
	case score >= 30:
		return "getting messy"
	case score >= 10:
		return "what are you doing bro"
	case score >= 0:
		return "fkn pain"
	case score >= -50:
		return "delete this shit"
	case score >= -100:
		return "is it that fucked up?"
	default:
		return "throw your computer away and delete the internet and never come back"
	}
}
