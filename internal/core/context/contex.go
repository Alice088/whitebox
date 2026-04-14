package context

import (
	"os"
	"path/filepath"
	"strings"
)

func mustBaseDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(home, ".whitebox")
}

var (
	BaseDir      = mustBaseDir()
	WorkspaceDir = filepath.Join(BaseDir, "workspace")
	ToolsDir     = filepath.Join(BaseDir, "context", "tools")
	SkillsDir    = filepath.Join(BaseDir, "context", "skills")
	MemoriesDir  = filepath.Join(BaseDir, "context", "memories")
	MindsDir     = filepath.Join(BaseDir, "context", "minds")
	SessionsDir  = filepath.Join(BaseDir, "context", "sessions")
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

	c.Minds, err = load(MindsDir)
	if err != nil {
		return err
	}

	c.Memories, err = load(MemoriesDir)
	if err != nil {
		return err
	}

	c.Skills, err = load(SkillsDir)
	if err != nil {
		return err
	}

	c.Tools, err = load(ToolsDir)
	if err != nil {
		return err
	}

	return nil
}
