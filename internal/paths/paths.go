package paths

import (
	"fmt"
	"path/filepath"
	"whitebox/pkg/meta"
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
	home := "./"
	return filepath.Join(home, fmt.Sprintf(".whitebox-%s", meta.AgentName))
}
