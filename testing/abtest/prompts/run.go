package prompts

import (
	"strings"
	"time"
	"whitebox/internal/core"
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

func (r *Runner) RunCase(c Case) Result {
	r.Logger.Info().Str("sys_prompt", trimSpace(messages.OutputLimit(c.Prompt, 5))).Msg("Mocking context")
	r.Engine.Context = abtest.NewContext(c.Prompt)

	var metrics Metrics
	start := time.Now()
	r.Logger.Info().Msg("Run engine")
	out, err := r.Engine.Run(c.Input, func(e core.Event) {
		metrics.EventCalls++

		if e.Type != "debug" {
			r.Logger.Info().Str("Type", e.Type).Int("Events", metrics.EventCalls).Msg("Event call")
		}

		switch e.Type {
		case "tool_call":
			metrics.ToolCalls++
		case "abtesting_loop_start":
			metrics.Steps++
		case "error":
			metrics.Errors++
		}
	})

	end := time.Since(start)
	metrics.Duration = end
	r.Logger.Info().Str("Duration", end.String()).Msg("End engine")

	return Result{
		Name:    c.Name,
		Output:  out,
		Error:   err,
		Metrics: metrics,
	}
}

func (r *Runner) RunBatch(cases []Case) []Result {
	var results []Result

	for _, c := range cases {
		r.Logger.Info().Str("name", c.Name).Msg("Run case")
		res := r.RunCase(c)
		results = append(results, res)
	}

	return results
}
