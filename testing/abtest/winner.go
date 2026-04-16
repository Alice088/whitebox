package abtest

import "sort"

func Compare(results []Result) []Result {
	out := make([]Result, len(results))
	copy(out, results)

	sort.Slice(out, func(i, j int) bool {
		s1, _ := Score(out[i].Metrics)
		s2, _ := Score(out[j].Metrics)
		return s1 > s2
	})

	return out
}

func PickWinner(results []Result) (Result, bool) {
	if len(results) == 0 {
		return Result{}, false
	}

	best := results[0]
	bestScore, _ := Score(best.Metrics)

	for _, r := range results[1:] {
		s, _ := Score(r.Metrics)
		if s > bestScore {
			best = r
			bestScore = s
		}
	}

	return best, true
}
