package context

import (
	"strings"
	"whitebox/internal/paths"
)

type Context interface {
	Prompt() string
	ClearMessages() error
	AddMessage(msg Message) error
	Collect() error
}

type Default struct {
	Sessions Session
	Tools    []Item
	Skills   []Item
	Memories []Item
	Minds    []Item
}

func New(session Session) Context {
	return &Default{
		Sessions: session,
	}
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

	for _, msg := range c.Sessions.Messages {
		builder.WriteString("\n")
		builder.WriteString(msg.Role)
		builder.WriteString(": ")
		builder.WriteString(msg.Content)
	}

	return builder.String()
}

func (c *Default) ClearMessages() error {
	c.Sessions.Messages = []Message{}
	return c.Sessions.SaveSession([]Message{})
}

func (c *Default) AddMessage(msg Message) error {
	c.Sessions.Messages = append(c.Sessions.Messages, msg)
	return c.Sessions.SaveSession(c.Sessions.Messages)
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
