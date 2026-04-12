package pipeline

import "github.com/henomis/langfuse-go/model"

type Observability struct {
	TraceID    string
	Generation *model.Generation
}
