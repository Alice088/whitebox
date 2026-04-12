package context

type Context struct {
	Messages []Item
	Tools    []Item
	Skills   []Item
	Memory   []Item
	Mind     []Item
	prompt   string
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
	if len(c.prompt) != 0 {
		return c.prompt
	}

	for _, item := range c.Mind {
		c.prompt += item.Content
	}

	for _, item := range c.Memory {
		c.prompt += item.Content
	}

	for _, item := range c.Skills {
		c.prompt += item.Content
	}

	for _, item := range c.Tools {
		c.prompt += item.Content
	}

	for _, item := range c.Messages {
		c.prompt += item.Content
	}

	return c.prompt
}

func (c *Context) Collect(opts CollectOpts) error {
	c.prompt = ""
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

	c.Messages, err = load(opts.MessagesPath)
	if err != nil {
		return err
	}

	return nil
}
