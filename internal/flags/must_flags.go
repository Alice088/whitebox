package flags

import "github.com/rs/zerolog"

func MustInput(logger zerolog.Logger) Input {
	input, err := ParseFlags()
	if err != nil {
		logger.Fatal().Err(err).Send()
	}
	return input
}
