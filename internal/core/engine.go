package core

import (
	"fmt"
	"strings"
	"time"
	"whitebox/internal/core/context"
	"whitebox/internal/core/llm"
	"whitebox/internal/core/tools"
	"whitebox/internal/langfuse"
	"whitebox/pkg/messages"
)

type Engine struct {
	LLM       llm.LLM
	Context   context.Context
	CallChain CallChain
}

func (e *Engine) Run(input string, emit func(Event)) (string, error) {
	state := &State{Input: input, Task: input}
	if w, ok := e.LLM.(*langfuse.LLMWrapper); ok {
		err := w.StartTrace(input)
		if err != nil {
			return "", fmt.Errorf("failed to start trace: %w", err)
		}
		defer w.EndTrace()
	}

	for i := 0; i < e.CallChain.Max; i++ {
		emit(Event{"abtesting_loop_start", fmt.Sprintf("loop start (i=%d)", i+1)})

		emit(Event{"debug", fmt.Sprintf("loop start (i=%d)", i+1)})
		emit(Event{"debug", fmt.Sprintf("call LLM (sys_prompt_len=%.2ft)", e.LLM.EstimateTokens(e.Context.Prompt()))})

		t := time.Now()
		out, err := e.LLM.Ask(state.Input, e.Context.Prompt())
		state.Output = out
		emit(Event{"debug", fmt.Sprintf("got response from LLM (%s)", time.Since(t).String())})

		if err != nil {
			emit(Event{"error", fmt.Sprintf("ASK_ERR: %v", err)})
			return "", err
		}
		emit(Event{"debug", fmt.Sprintf("LLM OUTPUT: [%s]", out)})

		if tc, ok := tools.IsToolCall(out); ok {
			result, err := e.LLM.Tool(tc)
			result = strings.TrimSpace(result)

			emit(Event{"abtesting_tool_call", tc})
			emit(Event{"tool_call", fmt.Sprintf("%s (%+v) \n - %s", tc.Tool, messages.StringArgs(tc.Arguments), messages.OutputLimit(result, 60))})
			state.History += fmt.Sprintf(`
						- Tool: %s
						  Args: %+v
						  Result: %s
						  Error: %v
			`, tc.Tool, tc.Arguments, result, err)

			if err != nil {
				emit(Event{Type: "error", Data: fmt.Sprintf("Call %s (%+v) \n - %s", tc.Tool, tc.Arguments, err.Error())})

				state.Input = fmt.Sprintf(`
					Tool "%s" executed with error.

					You must NOT use this pattern again
					Result:
					%s
					`, tc.Tool, err.Error())
			} else {
				state.Input = fmt.Sprintf(`
						You are solving a task step by step.
						
						Task:
						%s
						
						Previous actions:
						%s
						
						Last tool result:
						Tool: %s
						Result:
						%s
						
						Rules:
						
						Continue solving the task.
						
						If more actions are needed, call another tool.
						If task is complete, give final answer.
						
						Decide your next step.
					`, state.Task, state.History, tc.Tool, result) //todo это тоже в контекст перенести

			}
			continue
		}
		state.Input = out

		//if needsHuman(out) {
		//	emit(Event{"human_request", out})
		//
		//	answer := waitUser()
		//	emit(Event{"human_response", answer})
		//
		//	state.Input = answer
		//	continue
		//}
		emit(Event{"final", out})
		return out, nil
	}

	return "", fmt.Errorf("loop limit")
}
