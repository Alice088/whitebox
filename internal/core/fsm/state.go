package fsm

type State string

const (
	Idle     State = "idle"
	Analyze  State = "analyze"
	Plan     State = "plan"
	Act      State = "act"
	Observe  State = "observe"
	Finalize State = "finalize"
	Done     State = "done"
	Failed   State = "failed"
)

type Step struct {
	Title  string
	Done   bool
	Result string
}

type Machine struct {
	State State

	Iteration int
	MaxSteps  int

	CurrentStep int

	Memory WorkingMemory

	Errors []string
}

func New(maxSteps int) Machine {
	return Machine{
		State:    Idle,
		MaxSteps: maxSteps,
	}
}

func (m *Machine) Working() bool {
	return m.State != Done && m.State != Failed
}

func (m *Machine) Next() {
	switch m.State {
	case Idle:
		m.State = Analyze
	case Analyze:
		m.State = Plan
	case Plan:
		m.State = Act
	case Observe:
		m.State = Act
	case Finalize:
		m.State = Done
	}
}

func (m *Machine) MarkStepDone() {
	if m.CurrentStep >= len(m.Memory.Plan) {
		return
	}

	m.CurrentStep++
}
