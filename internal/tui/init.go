package tui

import (
	"fmt"
	"strings"
	"time"
	"whitebox/internal/core"
	"whitebox/internal/core/context"
	"whitebox/pkg/colors"

	"charm.land/bubbles/v2/cursor"
	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/atotto/clipboard"

	"github.com/alecthomas/chroma/v2/quick"
)

var inlineCodeStyle = lipgloss.NewStyle().
	Background(lipgloss.Color("236")).
	Foreground(lipgloss.Color("255")).
	Padding(0, 1)

type eventsMsg struct {
	events chan core.Event
}

type eventTickMsg struct {
	event core.Event
	ok    bool
}

type tuiModel struct {
	program     *tea.Program // 🔥 ВАЖНО
	chat        *Chat
	viewport    viewport.Model
	messages    []string
	textarea    textarea.Model
	senderStyle lipgloss.Style
	err         error
	loading     bool
	status      string
}

func renderInlineCode(s string) string {
	parts := strings.Split(s, "`")
	for i := 1; i < len(parts); i += 2 {
		parts[i] = inlineCodeStyle.Render(parts[i])
	}
	return strings.Join(parts, "")
}

func renderCodeBlocks(input string) string {
	lines := strings.Split(input, "\n")

	var out []string
	inCode := false
	lang := ""
	var code []string

	for _, line := range lines {
		if strings.HasPrefix(line, "```") {
			if inCode {
				var b strings.Builder
				err := quick.Highlight(&b, strings.Join(code, "\n"), lang, "terminal16m", "monokai")
				if err != nil {
					out = append(out, strings.Join(code, "\n"))
				} else {
					out = append(out, b.String())
				}
				code = nil
				inCode = false
				lang = ""
			} else {
				inCode = true
				lang = strings.TrimPrefix(line, "```")
			}
			continue
		}

		if inCode {
			code = append(code, line)
		} else {
			out = append(out, renderInlineCode(line))
		}
	}

	return strings.Join(out, "\n")
}

func initialModel(chat *Chat) tuiModel {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 1000
	ta.SetWidth(30)
	ta.SetHeight(3)

	s := ta.Styles()
	s.Focused.CursorLine = lipgloss.NewStyle()
	ta.SetStyles(s)

	ta.ShowLineNumbers = false
	ta.KeyMap.InsertNewline.SetEnabled(false)

	messages := []string{
		"Whitebox Chat Mode",
		fmt.Sprintf("Session ID: %s", chat.CoreEngine.Context.Session().ID),
		"Type '@exit' to quit, '@clear' to clear history",
	}

	vp := viewport.New(viewport.WithWidth(30), viewport.WithHeight(5))
	vp.SetContent(strings.Join(messages, "\n"))

	return tuiModel{
		chat:        chat,
		textarea:    ta,
		messages:    messages,
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
	}
}

func (m tuiModel) Init() tea.Cmd {
	return tea.Batch(textarea.Blink, tickCmd())
}

func askCmd(p *tea.Program, chat *Chat, input string) tea.Cmd {
	return func() tea.Msg {
		go func() {
			_, _ = chat.CoreEngine.Run(input, func(e core.Event) {
				p.Send(eventTickMsg{event: e, ok: true})
			})
			p.Send(eventTickMsg{ok: false})
		}()
		return nil
	}
}
func nextEventCmd(ch chan core.Event) tea.Cmd {
	return func() tea.Msg {
		e, ok := <-ch
		return eventTickMsg{event: e, ok: ok}
	}
}

func (m tuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.viewport.SetWidth(msg.Width)
		m.textarea.SetWidth(msg.Width)
		m.viewport.SetHeight(msg.Height - m.textarea.Height())
		m.viewport.SetContent(m.renderContent())
		m.viewport.GotoBottom()
		return m, nil

	case time.Time:
		if m.loading {
			m.status = m.chat.StatusEngine.NextAnimated()
			m.viewport.SetContent(m.renderContent())
		}
		return m, tickCmd()

	case tea.KeyPressMsg:
		switch msg.String() {

		case "ctrl+d", "esc":
			return m, tea.Quit
		case "ctrl+v":
			text, _ := clipboard.ReadAll()
			m.textarea.InsertString(text)

		case "enter":
			if m.loading {
				return m, nil
			}

			input := strings.TrimSpace(m.textarea.Value())
			if input == "" {
				return m, nil
			}

			m.messages = append(m.messages, m.senderStyle.Render("You: ")+input)
			m.chat.CoreEngine.Context.AddMessage(context.Message{
				Role:    "user",
				Content: input,
			})

			m.textarea.Reset()
			m.loading = true

			m.viewport.SetContent(m.renderContent())
			m.viewport.GotoBottom()

			return m, askCmd(m.program, m.chat, input)

		default:
			var cmd tea.Cmd
			m.textarea, cmd = m.textarea.Update(msg)

			var vpCmd tea.Cmd
			m.viewport, vpCmd = m.viewport.Update(msg)

			return m, tea.Batch(cmd, vpCmd)
		}

	case eventTickMsg:
		if msg.ok {
			switch msg.event.Type {
			case "debug":
				if m.chat.Debug {
					m.messages = append(m.messages, "DEBUG: "+fmt.Sprint(msg.event.Data))
				}

			case "llm_doing":
				m.messages = append(m.messages, msg.event.Data.(string))

			case "tool_call":
				m.messages = append(m.messages, colors.Dim("TOOLS: "+msg.event.Data.(string)))

			case "error":
				m.messages = append(m.messages, "ERROR: "+msg.event.Data.(string))

			case "final":
				rendered := renderCodeBlocks(msg.event.Data.(string))
				m.chat.CoreEngine.Context.AddMessage(context.Message{
					Role:    "assistant",
					Content: rendered,
				})
				m.messages = append(m.messages, "Assistant:\n"+rendered)
				m.loading = false
			}
		} else {
			m.loading = false
		}

		m.viewport.SetContent(m.renderContent())
		m.viewport.GotoBottom()

		return m, nil
	case cursor.BlinkMsg:
		var cmd tea.Cmd
		m.textarea, cmd = m.textarea.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m tuiModel) View() tea.View {
	divider := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render(strings.Repeat("─", m.viewport.Width()))

	input := lipgloss.NewStyle().
		MarginTop(1).
		Render(m.textarea.View())

	out := m.viewport.View() + "\n" + divider + "\n" + input

	if m.err != nil {
		out += "\nerror: " + m.err.Error()
	}

	v := tea.NewView(out)

	c := m.textarea.Cursor()
	if c != nil {
		c.Y += lipgloss.Height(m.viewport.View()) + 2 // 🔥 важно
	}

	v.Cursor = c
	v.AltScreen = true

	return v
}

func tickCmd() tea.Cmd {
	return tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
		return t
	})
}

func (m tuiModel) renderContent() string {
	content := strings.Join(m.messages, "\n")

	if m.loading && m.status != "" {
		content += "\n\n" + m.status
	}

	return lipgloss.NewStyle().
		Width(m.viewport.Width()).
		Render(content)
}
