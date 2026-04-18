package output

type Type string

const (
	FinalType Type = "final"
	ToolType  Type = "tool"
	PlanType  Type = "plan"
	AskType   Type = "ask"
)

type Answer[T any] struct {
	Type   Type
	Struct T
}

type Final struct {
	Type   string `json:"type"`
	Answer string `json:"answer"`
}

type Ask struct {
	Type string `json:"type"`
	Bool bool   `json:"bool"`
}

type Tool struct {
	Type      string            `json:"type"`
	ToolName  string            `json:"tool"`
	Arguments map[string]string `json:"arguments"`
}

type Plan struct {
	Type  string   `json:"type"`
	Steps []string `json:"steps"`
}
