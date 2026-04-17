package core

import (
	"fmt"
	"strings"
	"whitebox/internal/core/context"
	"whitebox/internal/core/embedded_prompts"
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

// todo добавть все побочныек промпты в langfuse
func (e *Engine) Run(input string, emit func(Event)) (string, error) {
	if w, ok := e.LLM.(*langfuse.LLMWrapper); ok {
		err := w.StartTrace(input)
		if err != nil {
			return "", fmt.Errorf("failed to start trace: %w", err)
		}
		defer w.EndTrace()
	}

	defer func() {
		if recv := recover(); recv != nil {
			emit(Event{"error", fmt.Sprintf("RUN_ERR: %+v", recv)})
		}
	}()

	machine := fsm.New(e.CallChain.Max)

	for machine.Working() {
		switch machine.State {

		case fsm.Idle:
			machine.Memory.Goal = input
			machine.Next()

		case fsm.Analyze:
			raw := strings.TrimSpace(machine.Memory.Goal)
			if raw == "" {
				machine.Errors = append(machine.Errors, "empty task")
				machine.State = fsm.Failed
				break
			}

			machine.Memory.Observations = append(
				machine.Memory.Observations,
				"task received",
			)
			machine.Next()

		case fsm.Plan:
			out, err := e.LLM.Ask(embedded_prompts.PlannerV1(machine.Memory.Goal), e.Context.Prompt())
			if err != nil {
				machine.Errors = append(machine.Errors, err.Error())
				machine.State = fsm.Failed
				break
			}

			answer, err := output.ToAnswer([]byte(out))
			if err != nil {
				machine.Errors = append(machine.Errors, "invalid planner json")
				machine.State = fsm.Failed
				break
			}

			if answer.Type != output.PlanType {
				machine.Errors = append(machine.Errors, "invalid planner json: not a plan")
				machine.State = fsm.Failed
				break
			}

			plan, _ := answer.Struct.(output.Plan)

			machine.Memory.Plan = plan.Steps
			machine.Memory.Observations = append(
				machine.Memory.Observations,
				"plan created",
			)
			machine.Next()

		case fsm.Act:
			prompt := buildPrompt(machine)

			out, err := e.LLM.Ask(prompt, e.Context.Prompt())
			if err != nil {
				machine.Errors = append(machine.Errors, err.Error())
				machine.State = fsm.Failed
				break
			}

			answer, err := output.ToAnswer([]byte(out))
			if err != nil {
				machine.Errors = append(machine.Errors, "invalid act json")
				machine.State = fsm.Failed
				break
			}

			switch answer.Type {

			case output.ToolType:
				toolCall, ok := answer.Struct.(output.Tool)
				if !ok {
					machine.Errors = append(machine.Errors, "invalid tool payload")
					machine.State = fsm.Failed
					break
				}

				result, toolErr := e.LLM.Tool(tools.ToolCall{
					Tool:      toolCall.ToolName,
					Arguments: toolCall.Arguments,
				})

				tr := fsm.ToolResult{
					Command: toolCall.ToolName,
					Output:  strings.TrimSpace(result),
					Success: toolErr == nil,
				}

				if toolErr != nil {
					tr.Error = toolErr.Error()
					machine.Errors = append(machine.Errors, toolErr.Error())
				}

				machine.Memory.ToolResults = append(
					machine.Memory.ToolResults,
					tr,
				)

				machine.Memory.LastAction = toolCall.ToolName
				machine.Memory.LastResult = tr.Output

				machine.State = fsm.Observe

			case output.FinalType:
				finalAnswer, ok := answer.Struct.(output.Final)
				if !ok {
					machine.Errors = append(machine.Errors, "invalid final payload")
					machine.State = fsm.Failed
					break
				}

				machine.Memory.LastResult = finalAnswer.Answer
				machine.State = fsm.Finalize

			case output.PlanType:
				plan, ok := answer.Struct.(output.Plan)
				if !ok {
					machine.Errors = append(machine.Errors, "invalid plan payload")
					machine.State = fsm.Failed
					break
				}

				machine.Memory.Plan = plan.Steps
				machine.Memory.Observations = append(
					machine.Memory.Observations,
					"plan updated",
				)

				machine.State = fsm.Act

			default:
				machine.Errors = append(machine.Errors, "unknown response type")
				machine.State = fsm.Failed
			}

		case fsm.Observe:
			machine.Memory.Observations = append(
				machine.Memory.Observations,
				machine.Memory.LastResult,
			)
			machine.State = fsm.Reflect

		case fsm.Reflect:
			machine.Memory.Attempts++

			if machine.Memory.Attempts >= machine.MaxSteps {
				machine.State = fsm.Failed
			} else {
				machine.State = fsm.Act
			}

		case fsm.Finalize:
			machine.State = fsm.Done
		}

		machine.Iteration++
	}

	if machine.State == fsm.Failed {
		return "", fmt.Errorf("failed: %s", strings.Join(machine.Errors, "; "))
	}

	return machine.Memory.LastResult, nil
}

func buildPrompt(m fsm.Machine) string {
	return fmt.Sprintf(`
		You are solving a task step by step.
		
		Goal:
		%s
		
		Plan:
		%s
		
		Observations:
		%s
		
		Last action:
		%s
		
		Last result:
		%s
		
		If needed, call a tool.
		If task is complete, return final answer only.
`,
		m.Memory.Goal,
		strings.Join(m.Memory.Plan, "\n"),
		strings.Join(m.Memory.Observations, "\n"),
		m.Memory.LastAction,
		m.Memory.LastResult,
	)
}
