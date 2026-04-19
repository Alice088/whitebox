package cfg

import (
	"path/filepath"
	"whitebox/internal/paths"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

func MustConfig(logger zerolog.Logger) Config {
	err := godotenv.Load(filepath.Join(paths.Base(), ".env"))
	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	var c Config
	err = env.Parse(&c)
	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	return c
}
