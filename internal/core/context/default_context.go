package context

import "github.com/rs/zerolog"

func NewMustDefault(msgs []Message, logger zerolog.Logger) *Context {
	context := Context{
		Messages: msgs,
	}
	err := context.Collect(CollectOpts{
		ToolsPath:  "./context/tools",
		MindPath:   "./context/mind",
		MemoryPath: "./context/memory",
		SkillsPath: "./context/skills",
	})

	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to collect context")
	}

	return &context
}
