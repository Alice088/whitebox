package status

import (
	"fmt"
	"math/rand"
	"time"
)

var statusWords = []string{
	"🔮  dream_oracle", "🧿  lucid_divining", "🌫️  night_whisper", "✨  summoning_dreams", "🌀  astral_channeling",
	"🪄  conjuring_visions", "📜  deciphering_dreams", "🧵  unraveling_sleep", "💭  deep_dreaming", "📡  signal_from_dream",
	"📢  echo_of_sleep", "🔍  seeking_visions", "🕳️  probing_dreamvoid", "📖  reading_dreams", "⚡  casting_visions",
	"⚙️  aligning_mind", "🎚️  tuning_sleep", "🌌  phasing_astral", "🔄  shifting_dream", "🧪  brewing_sleep",
	"⚒️  forging_visions", "🌅  awakening_dream", "📊  observing_subconscious", "🌀  folding_dreamspace", "👂  listening_within",
}

type StatusGenerator struct {
	lastWord    string
	rng         *rand.Rand
	colorIndex  int
	dotState    int
	currentWord string
}

func NewStatusGenerator() *StatusGenerator {
	sg := &StatusGenerator{
		rng:        rand.New(rand.NewSource(time.Now().UnixNano())),
		colorIndex: rand.New(rand.NewSource(time.Now().UnixNano())).Intn(216),
		dotState:   0,
	}
	sg.currentWord = sg.pickRandomWord()
	return sg
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

func (sg *StatusGenerator) pickRandomWord() string {
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
	return word
}

func (sg *StatusGenerator) NextAnimated() string {
	sg.colorIndex = (sg.colorIndex + 1) % 216
	colorCode := fmt.Sprintf("38;5;%d", sg.colorIndex+16)

	sg.dotState = (sg.dotState + 1) % 4
	dots := sg.getDots()

	return fmt.Sprintf("\033[%sm%s%s\033[0m", colorCode, sg.currentWord, dots)
}

func (sg *StatusGenerator) getDots() string {
	switch sg.dotState {
	case 0:
		return "."
	case 1:
		return ".."
	case 2:
		return "..."
	case 3:
		return ".."
	default:
		return "..."
	}
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
