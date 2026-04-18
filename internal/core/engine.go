package core

import (
	"fmt"
	"strings"
	"time"
	"whitebox/internal/core/context"
	"whitebox/internal/core/embedded_prompts"
	"whitebox/internal/core/fsm"
	"whitebox/internal/core/llm"
	"whitebox/internal/core/output"
	"whitebox/internal/core/tools"
	"whitebox/internal/langfuse"
	"whitebox/pkg/messages"
)

const (
	maxObservationItems = 30
	maxObservationSize  = 1880
	maxLastResultSize   = 5000
)

type Engine struct {
	LLM       llm.LLM
	Context   context.Context
	CallChain CallChain
}

func (e *Engine) Run(input string, emit func(Event)) (string, error) {
	runStarted := time.Now()

	emit(Event{"debug", "run:start"})
	emit(Event{"debug", fmt.Sprintf("model:%s", e.LLM.Model())})
	emit(Event{"debug", fmt.Sprintf("max_steps:%d", e.CallChain.Max)})

	if w, ok := e.LLM.(*langfuse.LLMWrapper); ok {
		err := w.StartTrace(input)
		if err != nil {
			return "", fmt.Errorf("failed to start trace: %w", err)
		}
		defer w.EndTrace()
	}

	defer func() {
		if recv := recover(); recv != nil {
			emit(Event{"error", fmt.Sprintf("panic:%+v", recv)})
		}
	}()

	machine := fsm.New(e.CallChain.Max)

	for machine.Working() {
		if machine.Iteration >= machine.MaxSteps {
			emit(Event{"debug", "max_steps_reached"})
			machine.State = fsm.Failed
		}

		emit(Event{
			"debug",
			fmt.Sprintf(
				"loop state=%s iter=%d/%d",
				machine.State,
				machine.Iteration,
				machine.MaxSteps,
			),
		})

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

			out, err := e.LLM.Ask(
				embedded_prompts.IsNeedPlanModeV1(raw),
				e.Context.Prompt(),
			)

			if err != nil {
				machine.Errors = append(machine.Errors, "failed classifier")
				machine.State = fsm.Failed
				break
			}

			answer, err := output.ToAnswer[output.Ask]([]byte(out))
			if err != nil {
				machine.Errors = append(machine.Errors, "invalid ask json")
				machine.State = fsm.Failed
				break
			}

			addObservation(&machine, "task received")

			if answer.Struct.Bool {
				machine.Next()
			} else {
				machine.State = fsm.Act
			}

		case fsm.Plan:
			var answer output.Answer[output.Plan]
			ok := false
			prompt := embedded_prompts.PlannerV1(machine.Memory.Goal)
			for i := 0; i < 3; i++ {
				out, err := e.LLM.Ask(prompt, e.Context.Prompt())

				if err != nil {
					machine.Errors = append(machine.Errors, err.Error())
					machine.State = fsm.Failed
					break
				}

				answer, err = output.ToAnswer[output.Plan]([]byte(out))
				if err != nil {
					if i == 2 {
						machine.Errors = append(machine.Errors, "invalid planner json")
						machine.State = fsm.Failed
					} else {
						emit(Event{
							"debug",
							fmt.Sprintf(
								"ERR_PLAN: %s ",
								"\nyour answer had mistake, fix previous error: "+fmt.Sprintf("invalid plan json: %s", err.Error()),
							),
						})

						prompt = embedded_prompts.PlannerV1(machine.Memory.Goal) +
							"\nyour answer had mistake, fix previous error: " + fmt.Sprintf("invalid plan json: %s", err.Error())
					}
					continue
				}
				ok = true
				break
			}

			if !ok {
				break
			}

			machine.Memory.Plan = normalizePlan(answer.Struct.Steps)
			machine.CurrentStep = 0
			machine.Next()

		case fsm.Act:
			machine.Iteration++

			out, err := e.LLM.Ask(
				buildPrompt(&machine),
				e.Context.Prompt(),
			)

			if err != nil {
				machine.Errors = append(machine.Errors, err.Error())
				machine.State = fsm.Failed
				break
			}

			if strings.TrimSpace(out) == "" {
				machine.Memory.LastResult = "stopped: empty model response"
				machine.State = fsm.Finalize
				break
			}

			answer, err := output.ToAnswer[any]([]byte(out))
			if err != nil {
				machine.Errors = append(machine.Errors, "invalid act json")
				machine.State = fsm.Failed
				break
			}

			switch answer.Type {

			case output.ToolType:
				toolCall := answer.Struct.(output.Tool)

				emit(Event{
					"tool_call",
					fmt.Sprintf(
						"%s %+v",
						toolCall.ToolName,
						messages.StringArgs(toolCall.Arguments),
					),
				})

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

				machine.Memory.ToolResults = append(machine.Memory.ToolResults, tr)

				if tr.Success {
					machine.MarkStepDone()
				}

				machine.Memory.LastAction = toolCall.ToolName
				machine.Memory.LastResult = summarizeResult(
					toolCall.ToolName,
					tr.Output,
				)

				machine.State = fsm.Observe

			case output.FinalType:
				machine.CurrentStep = len(machine.Memory.Plan)
				machine.Memory.LastResult =
					answer.Struct.(output.Final).Answer
				machine.State = fsm.Finalize

			case output.PlanType:
				addObservation(&machine, "replanning")
				machine.Errors = append(
					machine.Errors,
					"replanning don't allow in act",
				)

				machine.State = fsm.Failed

			default:
				machine.Errors = append(
					machine.Errors,
					"unknown response type",
				)
				machine.State = fsm.Failed
			}

		case fsm.Observe:
			addObservation(
				&machine,
				machine.Memory.LastAction+" => "+machine.Memory.LastResult,
			)

			if machine.CurrentStep >= len(machine.Memory.Plan) {
				machine.State = fsm.Finalize
				break
			}

			machine.Next()

		case fsm.Finalize:
			emit(Event{
				"final",
				machine.Memory.LastResult,
			})

			machine.State = fsm.Done
		}
	}

	if machine.State == fsm.Failed {
		errText := fmt.Sprintf(
			"failed: %s",
			strings.Join(machine.Errors, "; "),
		)

		return "", fmt.Errorf(errText)
	}

	emit(Event{
		"debug",
		fmt.Sprintf(
			"run:done total_ms=%d",
			time.Since(runStarted).Milliseconds(),
		),
	})

	return machine.Memory.LastResult, nil
}

func buildPrompt(m *fsm.Machine) string {
	action := "finish task"

	if m.CurrentStep < len(m.Memory.Plan) {
		action = m.Memory.Plan[m.CurrentStep]
	}

	return fmt.Sprintf(`
IMPORTANT:
You are currently in EXECUTION MODE.

Allowed response types ONLY:
- tool
- final

Do NOT return:
- plan
- ask

Return one JSON object only.

Goal:
%s

Plan:
%s

Current step:
%d/%d

Next action:
%s

Observations:
%s

Last action:
%s

Last result:
%s

Rules:
- Follow current plan
- Execute next unfinished step
- Do not repeat completed steps
- Prefer tool call when action is needed
- If task complete return final
- DO NOT return final if git status --short is not empty
- If any file is:
	- staged
	- untracked
	- modified
- You MUST create commits until repository is clean.
If there are untracked files (??):
→ you MUST run git add before commit
`,
		m.Memory.Goal,
		messages.FlatArr(m.Memory.Plan),
		m.CurrentStep+1, len(m.Memory.Plan),
		action,
		strings.Join(m.Memory.Observations, "\n"),
		m.Memory.LastAction,
		compact(m.Memory.LastResult, maxLastResultSize),
	)
}

func addObservation(m *fsm.Machine, value string) {
	m.Memory.Observations = append(
		m.Memory.Observations,
		compact(value, maxObservationSize),
	)

	if len(m.Memory.Observations) > maxObservationItems {
		m.Memory.Observations =
			m.Memory.Observations[len(m.Memory.Observations)-maxObservationItems:]
	}
}

func summarizeResult(action, result string) string {
	raw := strings.TrimSpace(result)

	if raw == "" {
		return "ok"
	}

	lines := strings.Count(raw, "\n") + 1

	return fmt.Sprintf(
		"%s output (%d chars, %d lines): %s",
		action,
		len(raw),
		lines,
		compact(raw, maxLastResultSize),
	)
}

func compact(s string, limit int) string {
	s = strings.TrimSpace(s)

	if len(s) <= limit {
		return s
	}

	if limit < 4 {
		return s[:limit]
	}

	return s[:limit-3] + "..."
}

func normalizePlan(steps []string) []string {
	out := make([]string, 0, len(steps))

	for i, step := range steps {
		s := strings.TrimSpace(step)
		if s == "" {
			continue
		}
		s = fmt.Sprintf("(%d) %s", i+1, s)
		out = append(out, s)
	}

	if len(out) == 0 {
		return []string{"finish task"}
	}

	return out
}
