package context

func NewDefault() (Context, error) {
	context := Context{
		Messages: []Message{},
	}
	return context, context.Collect(CollectOpts{
		ToolsPath:  "./context/tools",
		MindPath:   "./context/mind",
		MemoryPath: "./context/memory",
		SkillsPath: "./context/skills",
	})
}
