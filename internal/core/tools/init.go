package tools

func init() {
	Tools["read_file"] = ReadFile
	Tools["memory"] = WriteFile
	Tools["write_file"] = WriteFile
	Tools["bash"] = Bash
}

type Tool func(arguments map[string]string) (string, error)

var Tools = make(map[string]Tool)

type ToolCall struct {
	Tool      string            `json:"Tool"`
	Arguments map[string]string `json:"arguments"`
}
