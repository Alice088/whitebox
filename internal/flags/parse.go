package flags

import (
	"flag"
)

type Input struct {
	SessionID string
	Debug     bool
}

func ParseFlags() (Input, error) {
	debug := flag.Bool("debug", false, "debug mode without tui")

	sessionID := flag.String("session", "", "session ID for persistent chat")
	flag.Parse()

	return Input{
		SessionID: *sessionID,
		Debug:     *debug,
	}, nil
}
