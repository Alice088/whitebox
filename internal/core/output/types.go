package output

type Type string

const (
	FinalType Type = "final"
	ToolType  Type = "tool"
	PlanType  Type = "plan"
)

type Answer struct {
	Type   Type
	Struct any
}

type Final struct {
	Type   string `json:"type"`
	Answer string `json:"answer"`
}

type Tool struct {
	Type      string            `json:"type"`
	ToolName  string            `json:"tool_name"`
	Arguments map[string]string `json:"arguments"`
}

type Plan struct {
	Type  string   `json:"type"`
	Steps []string `json:"steps"`
}
