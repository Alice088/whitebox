package prompts

import (
	"strings"
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

	r.Logger.Info().Msg("Run engine")
	out, err := r.Engine.Run(c.Input, func(e core.Event) {})

	return Result{
		Name:   c.Name,
		Output: out,
		Error:  err,
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
