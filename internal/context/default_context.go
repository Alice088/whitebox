package context

func NewDefault() (Context, error) {
	context := Context{}
	return context, context.Collect(CollectOpts{
		ToolsPath:  "./context/tools",
		MindPath:   "./context/mind",
		MemoryPath: "./context/memory",
		//MessagesPath: "nope",
		SkillsPath: "./context/skills",
	})
}
