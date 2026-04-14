package tui

import (
	"fmt"
	"whitebox/internal/core"
	"whitebox/internal/tui/status"

	tea "charm.land/bubbletea/v2"
)

type Chat struct {
	CoreEngine   core.Engine
	StatusEngine *status.StatusEngine
	Debug        bool
}

func New(engine core.Engine, debug bool) Chat {
	return Chat{
		CoreEngine:   engine,
		Debug:        debug,
		StatusEngine: status.NewStatusEngine(),
	}
}

func (chat *Chat) Run() {
	m := initialModel(chat)
	p := tea.NewProgram(&m)
	m.program = p

	if _, err := p.Run(); err != nil {
		fmt.Println("failed:", err)
	}
}
