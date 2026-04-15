package abtest

import (
	"whitebox/internal/core/context"
)

type Context struct {
	SystemPrompt string
	Tools        []context.Item
	Skills       []context.Item
	Memories     []context.Item
	Minds        []context.Item
}

func NewContext(systemPrompt string) context.Context {
	return &Context{
		SystemPrompt: systemPrompt,
	}
}

func (c *Context) Prompt() string {
	return c.SystemPrompt
}

func (c *Context) ClearMessages() error {
	return nil
}

func (c *Context) AddMessage(msg context.Message) error {
	return nil
}

func (c *Context) Collect() error {
	return nil
}
