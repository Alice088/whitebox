package flags

import (
	"flag"
)

type Input struct {
	SessionID string
	Debug     bool
	Headless  bool
	Msg       string
}

func ParseFlags() (Input, error) {
	debug := flag.Bool("debug", false, "debug")
	headless := flag.Bool("headless", false, "engine without tui")
	msg := flag.String("msg", "", "msg to llm")

	sessionID := flag.String("session", "", "session ID for persistent chat")
	flag.Parse()

	return Input{
		SessionID: *sessionID,
		Debug:     *debug,
		Headless:  *headless,
		Msg:       *msg,
	}, nil
}
