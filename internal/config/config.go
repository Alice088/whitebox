package config

type Config struct {
	LLM LLM
}

type LLM struct {
	ApiKey string `env:"LLM_API_KEY,required"`
}
