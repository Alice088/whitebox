package context

import (
	"strings"
	"whitebox/internal/paths"
)

type Context struct {
	Sessions Sessions
	Tools    []Item
	Skills   []Item
	Memories []Item
	Minds    []Item
}

type Item struct {
	Source  string
	Content string
}

func (c *Context) Prompt() string {
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

func (c *Context) ClearMessages() error {
	c.Sessions.Messages = []Message{}
	return c.Sessions.SaveSession([]Message{})
}

func (c *Context) AddMessage(msg Message) error {
	c.Sessions.Messages = append(c.Sessions.Messages, msg)
	return c.Sessions.SaveSession(c.Sessions.Messages)
}

func (c *Context) Collect() error {
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
