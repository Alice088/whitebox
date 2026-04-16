package prompts

import (
	"strings"
	"time"
	"whitebox/internal/core"
	"whitebox/internal/core/tools"
	"whitebox/pkg/maps"
	"whitebox/pkg/messages"
	"whitebox/testing/abtest"

	"github.com/rs/zerolog"
)

type Runner struct {
	Engine *core.Engine
	Logger *zerolog.Logger
}

func trimSpace(str string) string {
	return strings.ReplaceAll(strings.TrimSpace(str), "\n", " ")
}

func (r *Runner) RunCase(c abtest.Case) abtest.Result {
	r.Logger.Info().Str("sys_prompt", trimSpace(messages.OutputLimit(c.Prompt, 100))).Msg("Mocking context")
	r.Engine.Context = abtest.NewContext(c.Prompt)

	var metrics abtest.Metrics
	metrics.ToolsCalls = make(map[string]int)
	start := time.Now()
	r.Logger.Info().Msg("Run engine")
	out, err := r.Engine.Run(c.Input, func(e core.Event) {
		metrics.EventCalls++

		if e.Type != "debug" {
			r.Logger.Info().Str("Type", e.Type).Int("Events", metrics.EventCalls).Msg("Event call")
		}

		switch e.Type {
		case "abtesting_tool_call":
			tc := e.Data.(tools.ToolCall)

			metrics.ToolsCallsHistory = append(metrics.ToolsCallsHistory, tc)

			if maps.Exists(metrics.ToolsCalls, tc.Tool) {
				metrics.ToolsCalls[tc.Tool]++
			} else {
				metrics.ToolsCalls[tc.Tool] = 1
			}
		case "abtesting_loop_start":
			metrics.Steps++
		case "error":
			metrics.Errors++
		}
	})

	end := time.Since(start)
	metrics.Duration = end
	r.Logger.Info().Str("Duration", end.String()).Msg("End engine")

	return abtest.Result{
		Name:    c.Name,
		Output:  out,
		Error:   err,
		Metrics: metrics,
	}
}

func (r *Runner) RunBatch(cases []abtest.Case) []abtest.Result {
	var results []abtest.Result

	for _, c := range cases {
		r.Logger.Info().Str("name", c.Name).Msg("Run case")
		res := r.RunCase(c)
		results = append(results, res)
	}

	return results
}
