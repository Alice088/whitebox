package core

import (
	"fmt"
	"time"
	"whitebox/internal/core/context"
	"whitebox/internal/core/llm"
	"whitebox/internal/core/tools"
)

type Engine struct {
	LLM       llm.LLM
	Context   *context.Context
	CallChain CallChain
}

func (e *Engine) Run(input string, emit func(Event)) (string, error) {
	state := &State{Input: input}

	for i := 0; i < e.CallChain.Max; i++ {
		emit(Event{"debug", fmt.Sprintf("loop start (i=%d)", i+1)})
		emit(Event{"debug", fmt.Sprintf("call LLM (sys_prompt_len=%.2ft)", e.LLM.EstimateTokens(e.Context.Prompt()))})

		t := time.Now()
		out, err := e.LLM.Ask(state.Input, e.Context.Prompt())
		state.Output = out
		emit(Event{"debug", fmt.Sprintf("got response from LLM (%s)", time.Since(t).String())})

		if err != nil {
			emit(Event{"error", fmt.Sprintf("Error!: [%v]", err)})
			return "", err
		}
		emit(Event{"debug", fmt.Sprintf("LLM OUTPUT: [%s]", out)})

		if tc, ok := tools.TryParseToolCall(out); ok {
			result, err := tools.Execute(tc)
			emit(Event{"tool_call", fmt.Sprintf("%s (%+v) \n - %s", tc.Tool, tc.Arguments, result)})
			if err != nil {
				state.Input = fmt.Sprintf("Tool result(%s; args: %+v): %s", tc.Tool, tc.Arguments, err.Error())
			} else {
				state.Input = fmt.Sprintf("Tool result(%s; args: %+v): %s", tc.Tool, tc.Arguments, result)
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
