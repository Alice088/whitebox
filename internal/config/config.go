package config

type Config struct {
	LLM           LLM
	Observability Observability
	Session       Session
	CallChain     CallChain
}

type LLM struct {
	ApiKey  string `env:"LLM_API_KEY,required"`
	Provide string `env:"LLM_PROVIDER,required"`
	Model   string `env:"LLM_MODEL,required"`
}

type Observability struct {
	LangFuse LangFuse
}

type LangFuse struct {
	Enabled bool `env:"LANGFUSE_ENABLED,required"`
}

type Session struct {
	MaxMessages int `env:"SESSION_MAX_MESSAGES,required"`
}
type CallChain struct {
	Max int `env:"CALL_CHAIN_MAX,required"`
}
