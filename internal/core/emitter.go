package core

type Emitter interface {
	Emit(e Event)
}

type Event struct {
	Type string
	Data any
}
