package fsm

type State string

const (
	Idle     State = "idle"
	DoOne    State = "do_one"
	DoSecond State = "do_second"
	Done     State = "done"
	Final    State = "final"
	Fail     State = "fail"
)

type Step struct {
	Title  string
	Done   bool
	Result string
}

type Machine struct {
	State  State
	Memory WorkingMemory
}

func New(maxSteps int) Machine {
	return Machine{
		State: Idle,
	}
}

func (m *Machine) Working() bool {
	return m.State != Done && m.State != Fail
}

func (m *Machine) Next() {
	switch m.State {
	case Idle:
		m.State = DoOne
	case DoOne:
		m.State = DoSecond
	case DoSecond:
		m.State = Final
	}
}
