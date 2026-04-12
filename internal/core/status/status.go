package status

import (
	"fmt"
	"math/rand"
	"time"
)

var statusWords = []string{
	"oracle", "divining", "whispering", "summoning", "channeling",
	"conjuring", "deciphering", "unraveling", "dreaming", "glitching",
	"echoing", "seeking", "probing_void", "reading_runes", "casting",
	"aligning", "tuning", "phasing", "shifting", "brewing",
	"forging", "awakening", "observing_entropy", "folding_space", "listening",
}

type StatusGenerator struct {
	lastWord string
	rng      *rand.Rand
}

func NewStatusGenerator() *StatusGenerator {
	return &StatusGenerator{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (sg *StatusGenerator) Next() string {
	var available []string
	for _, word := range statusWords {
		if word != sg.lastWord {
			available = append(available, word)
		}
	}

	if len(available) == 0 {
		available = statusWords
	}

	word := available[sg.rng.Intn(len(available))]
	sg.lastWord = word

	colorCode := sg.randomColor()

	if sg.rng.Float32() < 0.3 {
		style := sg.rng.Intn(2) + 1
		return fmt.Sprintf("\033[%d;%sm%s\033[0m", style, colorCode, word)
	}

	return fmt.Sprintf("\033[%sm%s\033[0m", colorCode, word)
}

func (sg *StatusGenerator) randomColor() string {
	switch sg.rng.Intn(3) {
	case 0:
		return fmt.Sprintf("%d", sg.rng.Intn(7)+31)
	case 1:
		return fmt.Sprintf("%d", sg.rng.Intn(6)+91)
	case 2:
		n := sg.rng.Intn(240) + 16
		return fmt.Sprintf("38;5;%d", n)
	}
	return "37"
}
