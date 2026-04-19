package output

type Type string

const (
	FinalType Type = "final"
	ToolType  Type = "tool"
)

type Answer[T any] struct {
	Type   Type
	Struct T
}

type Final struct {
	Type   string `json:"type"`
	Answer string `json:"answer"`
}

type Tool struct {
	Type      string            `json:"type"`
	ToolName  string            `json:"tool"`
	Arguments map[string]string `json:"arguments"`
}
