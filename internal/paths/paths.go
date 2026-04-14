package paths

import (
	"os"
	"path/filepath"
)

var (
	BaseDir      = Base()
	WorkspaceDir = filepath.Join(BaseDir, "workspace")
	ToolsDir     = filepath.Join(BaseDir, "context", "tools")
	SkillsDir    = filepath.Join(BaseDir, "context", "skills")
	MemoriesDir  = filepath.Join(BaseDir, "context", "memories")
	MindsDir     = filepath.Join(BaseDir, "context", "minds")
	SessionsDir  = filepath.Join(BaseDir, "context", "sessions")
	CommandsDir  = filepath.Join(BaseDir, "commands")
)

func Base() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic("failed to load home dir")
	}

	return filepath.Join(home, ".whitebox")
}
