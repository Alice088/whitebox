package core

import (
	"whitebox/internal/core/context"
	"whitebox/internal/core/fsm"
	"whitebox/internal/core/llm"
	"whitebox/internal/core/output"
	"whitebox/internal/core/tools"
	"whitebox/internal/langfuse"
)

type Engine struct {
	LLM       llm.LLM
	Context   context.Context
	CallChain CallChain
}

func (e *Engine) Run(input string, emit func(Event)) (string, error) {
	if w, ok := e.LLM.(*langfuse.LLMWrapper); ok {
		if err := w.StartTrace(input); err != nil {
			return "", err
		}
		defer w.EndTrace()
	}

	m := fsm.New(e.CallChain.Max)
	m.Memory.Goal = input

	for m.Working() {
		switch m.State {

		case fsm.Idle:
			m.Next()

		case fsm.DoOne:
			out, _ := e.LLM.Ask(input, e.Context.Prompt())

			answer, _ := output.ToAnswer[any]([]byte(out))

			switch answer.Type {

			case output.ToolType:
				toolCall := answer.Struct.(output.Tool)

				result, _ := e.LLM.Tool(tools.ToolCall{
					Tool:      toolCall.ToolName,
					Arguments: toolCall.Arguments,
				})

				m.Memory.LastResult = result
				m.Next()

			case output.FinalType:
				m.Memory.LastResult = answer.Struct.(output.Final).Answer
				m.State = fsm.Final
			}

		case fsm.DoSecond:
			out, _ := e.LLM.Ask(input, e.Context.Prompt())

			answer, _ := output.ToAnswer[any]([]byte(out))

			switch answer.Type {

			case output.ToolType:
				toolCall := answer.Struct.(output.Tool)

				result, _ := e.LLM.Tool(tools.ToolCall{
					Tool:      toolCall.ToolName,
					Arguments: toolCall.Arguments,
				})

				m.Memory.LastResult = result
				m.Next()

			case output.FinalType:
				m.Memory.LastResult = answer.Struct.(output.Final).Answer
				m.State = fsm.Final
			}

		case fsm.Final:
			emit(Event{"final", m.Memory.LastResult})
			return m.Memory.LastResult, nil
		}
	}

	return m.Memory.LastResult, nil
}
