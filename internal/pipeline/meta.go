package pipeline

import "github.com/henomis/langfuse-go/model"

type Meta struct {
	TraceID    string
	Generation *model.Generation
}
