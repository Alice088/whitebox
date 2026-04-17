package fsm

type WorkingMemory struct {
	Goal string
	Plan []string

	Observations []string
	ToolResults  []ToolResult

	LastAction string
	LastResult string

	Attempts int
}

type ToolResult struct {
	Command  string
	Output   string
	Error    string
	Success  bool
	Duration int64
}
