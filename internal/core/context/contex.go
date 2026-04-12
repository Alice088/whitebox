package context

import (
	"strings"
)

type Context struct {
	Messages []Message
	Tools    []Item
	Skills   []Item
	Memory   []Item
	Mind     []Item

	staticPrompt   string
	messagesPrompt string
	staticValid    bool
	messagesValid  bool
}

type Item struct {
	Source  string
	Content string
}

type CollectOpts struct {
	MemoryPath   string
	SkillsPath   string
	ToolsPath    string
	MessagesPath string
	MindPath     string
}

func (c *Context) Prompt() string {
	c.buildStatic()
	c.buildMessages()
	return c.staticPrompt + c.messagesPrompt
}

func (c *Context) buildStatic() {
	if c.staticValid {
		return
	}

	var builder strings.Builder

	for _, item := range c.Mind {
		builder.WriteString(item.Content)
	}

	for _, item := range c.Memory {
		builder.WriteString(item.Content)
	}

	for _, item := range c.Skills {
		builder.WriteString(item.Content)
	}

	for _, item := range c.Tools {
		builder.WriteString(item.Content)
	}

	c.staticPrompt = builder.String()
	c.staticValid = true
}

func (c *Context) buildMessages() {
	if c.messagesValid {
		return
	}

	var builder strings.Builder

	for _, msg := range c.Messages {
		builder.WriteString("\n")
		builder.WriteString(msg.Role)
		builder.WriteString(": ")
		builder.WriteString(msg.Content)
	}

	c.messagesPrompt = builder.String()
	c.messagesValid = true
}

func (c *Context) AddMessage(msg Message) {
	c.Messages = append(c.Messages, msg)
	c.messagesValid = false
}

func (c *Context) ClearMessages() {
	c.Messages = []Message{}
	c.messagesPrompt = ""
	c.messagesValid = false
}

func (c *Context) TrimMessages(max int) {
	if len(c.Messages) > max {
		c.Messages = c.Messages[len(c.Messages)-max:]
		c.messagesValid = false
	}
}

func (c *Context) SetMind(items []Item) {
	c.Mind = items
	c.staticValid = false
}

func (c *Context) SetMemory(items []Item) {
	c.Memory = items
	c.staticValid = false
}

func (c *Context) SetSkills(items []Item) {
	c.Skills = items
	c.staticValid = false
}

func (c *Context) SetTools(items []Item) {
	c.Tools = items
	c.staticValid = false
}

func (c *Context) Collect(opts CollectOpts) error {
	c.staticValid = false
	c.messagesValid = false
	var err error

	c.Mind, err = load(opts.MindPath)
	if err != nil {
		return err
	}

	c.Memory, err = load(opts.MemoryPath)
	if err != nil {
		return err
	}

	c.Skills, err = load(opts.SkillsPath)
	if err != nil {
		return err
	}

	c.Tools, err = load(opts.ToolsPath)
	if err != nil {
		return err
	}

	return nil
}
