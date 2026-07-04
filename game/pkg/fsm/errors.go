package fsm

type MachineError string

func (e MachineError) Error() string {
	return string(e)
}

const (
	// ErrStateNotFound operation state not found inside machine.
	ErrStateNotFound = MachineError("state not found")

	// ErrMachineNotRunning operation not allowed for stopped machine.
	ErrMachineNotRunning = MachineError("machine not running")

	// ErrUnknownEvent event not found in transitions table.
	ErrUnknownEvent = MachineError("unknown event")

	// ErrNoStates no states registered in the machine.
	ErrNoStates = MachineError("no states found")

	// ErrTransitionLeftStateNotFound transition source state is not registered.
	ErrTransitionLeftStateNotFound = MachineError("transition left state not found")

	// ErrTransitionRightStateNotFound transition target state is not registered.
	ErrTransitionRightStateNotFound = MachineError("transition right state not found")

	// ErrFinalStateNotFound final state is not registered.
	ErrFinalStateNotFound = MachineError("final state not found")

	// ErrInitialStateNotInitialized initial state was not set on the builder.
	ErrInitialStateNotInitialized = MachineError("initial state not initialized")

	// ErrInitialStateNotRegistered initial state is not registered.
	ErrInitialStateNotRegistered = MachineError("initial state not registered")
)
