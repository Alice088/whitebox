package context

import (
	"strings"
	"whitebox/internal/paths"
)

type Context interface {
	Prompt() string
	Session() Session
	ClearMessages() error
	AddMessage(msg Message) error
	Collect() error
}

type Default struct {
	S               Session
	Tools           []Item
	Skills          []Item
	Memories        []Item
	Minds           []Item
	EmbeddedPrompts []string
}

func New(session Session, embeddedPrompts ...string) Context {
	d := &Default{
		S: session,
	}

	for _, prompt := range embeddedPrompts {
		d.EmbeddedPrompts = append(d.EmbeddedPrompts, prompt)
	}

	return d
}

type Item struct {
	Source  string
	Content string
}

func (c *Default) Prompt() string {
	var builder strings.Builder

	for _, item := range c.Minds {
		builder.WriteString(item.Content)
	}

	for _, item := range c.Memories {
		builder.WriteString(item.Content)
	}

	for _, item := range c.Skills {
		builder.WriteString(item.Content)
	}

	for _, item := range c.Tools {
		builder.WriteString(item.Content)
	}
	for _, prompt := range c.EmbeddedPrompts {
		builder.WriteString(prompt)
	}

	for _, msg := range c.S.Messages {
		builder.WriteString("\n")
		builder.WriteString(msg.Role)
		builder.WriteString(": ")
		builder.WriteString(msg.Content)
	}

	return builder.String()
}

func (c *Default) ClearMessages() error {
	c.S.Messages = []Message{}
	return c.S.SaveSession([]Message{})
}

func (c *Default) AddMessage(msg Message) error {
	c.S.Messages = append(c.S.Messages, msg)
	return c.S.SaveSession(c.S.Messages)
}

func (c *Default) Collect() error {
	var err error

	c.Minds, err = load(paths.MindsDir)
	if err != nil {
		return err
	}

	c.Memories, err = load(paths.MemoriesDir)
	if err != nil {
		return err
	}

	c.Skills, err = load(paths.SkillsDir)
	if err != nil {
		return err
	}

	c.Tools, err = load(paths.ToolsDir)
	if err != nil {
		return err
	}

	return nil
}

func (c *Default) Session() Session {
	return c.S
}
