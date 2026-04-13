package tui

import (
	"context"
	"fmt"
	"strings"
	"time"
	syscontext "whitebox/internal/core/context"

	"charm.land/bubbles/v2/cursor"
	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/alecthomas/chroma/v2/quick"
)

var inlineCodeStyle = lipgloss.NewStyle().
	Background(lipgloss.Color("236")).
	Foreground(lipgloss.Color("255")).
	Padding(0, 1)

type answerMsg struct {
	text string
	err  error
}

type tuiModel struct {
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

func initialModel(chat *Chat, sessionID string) tuiModel {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.SetVirtualCursor(false)
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
		fmt.Sprintf("Session ID: %s", sessionID),
		"Type '@exit' to quit, '@clear' to clear history",
	}

	vp := viewport.New(viewport.WithWidth(30), viewport.WithHeight(5))

	vp.SetContent(
		lipgloss.NewStyle().
			Width(30).
			Render(strings.Join(messages, "\n")),
	)

	return tuiModel{
		chat:        chat,
		textarea:    ta,
		messages:    messages,
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
	}
}

func (m tuiModel) Init() tea.Cmd {
	return tea.Batch(
		textarea.Blink,
		tickCmd(),
	)
}

func askCmd(chat *Chat, input string) tea.Cmd {
	return func() tea.Msg {
		out, err := chat.ask(context.Background(), input)
		return answerMsg{text: out, err: err}
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
			m.status = m.chat.statusEngine.NextAnimated()
			m.viewport.SetContent(m.renderContent())
		}
		return m, tickCmd()

	case tea.KeyPressMsg:
		switch msg.String() {

		case "ctrl+c", "esc":
			return m, tea.Quit

		case "enter":
			if m.loading {
				return m, nil
			}

			input := strings.TrimSpace(m.textarea.Value())
			if input == "" {
				return m, nil
			}

			if input == "@exit" {
				return m, tea.Quit
			}

			if input == "@clear" {
				m.chat.Context.ClearMessages()
				m.messages = m.messages[:3]
				m.viewport.SetContent(m.renderContent())
				m.textarea.Reset()
				return m, nil
			}

			m.chat.Context.AddMessage(syscontext.Message{
				Role:    "user",
				Content: input,
			})
			m.chat.Context.TrimMessages(m.chat.Session.MaxMessages)

			m.messages = append(m.messages,
				m.senderStyle.Render("You: ")+input,
			)

			m.textarea.Reset()
			m.loading = true
			m.status = m.chat.statusEngine.NextAnimated()

			m.viewport.SetContent(m.renderContent())
			m.viewport.GotoBottom()

			return m, tea.Batch(
				askCmd(m.chat, input),
				tickCmd(),
			)

		default:
			var cmd tea.Cmd
			m.textarea, cmd = m.textarea.Update(msg)

			var vpCmd tea.Cmd
			m.viewport, vpCmd = m.viewport.Update(msg)

			return m, tea.Batch(cmd, vpCmd)
		}

	case answerMsg:
		m.loading = false
		m.status = ""

		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}

		rendered := renderCodeBlocks(msg.text)
		m.messages = append(m.messages, "Assistant:\n"+rendered)

		m.viewport.SetContent(m.renderContent())
		m.viewport.GotoBottom()

		m.chat.Context.AddMessage(syscontext.Message{
			Role:    "assistant",
			Content: msg.text,
		})
		m.chat.Context.TrimMessages(m.chat.Session.MaxMessages)

		if len(m.chat.Context.Messages) > 0 {
			if err := m.chat.Session.SaveSession(m.chat.Context.Messages); err != nil {
				m.err = err
			}
		}

		return m, nil

	case cursor.BlinkMsg:
		var cmd tea.Cmd
		m.textarea, cmd = m.textarea.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m tuiModel) View() tea.View {
	viewportView := m.viewport.View()

	out := viewportView + "\n" + m.textarea.View()

	if m.err != nil {
		out += "\nerror: " + m.err.Error()
	}

	v := tea.NewView(out)

	c := m.textarea.Cursor()
	if c != nil {
		c.Y += lipgloss.Height(viewportView)
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
