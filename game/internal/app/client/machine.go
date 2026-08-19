package client

import (
	"engine/window"
	"game/internal/gamestate/gamemachine"
	"game/internal/rendering"
)

type MachineModules struct {
	Window    *window.Window
	Rendering rendering.Rendering
}

func NewMachine(mm MachineModules) gamemachine.Machine {
	return gamemachine.NewMachineBuilder().
		InitialState(gamemachine.PlayingStateIdentifier). // todo change to InitialState
		BuildState(gamemachine.NewInitialState()).
		Build().
		BuildState(gamemachine.NewResourcesLoadingState()).
		Transition(gamemachine.ResourcesLoadedEvent, gamemachine.PlayingStateIdentifier).
		Build().
		BuildState(gamemachine.NewPlayingState(mm.Rendering)).
		Build().
		RegisterState(gamemachine.NewStoppedState(), make(gamemachine.Transitions)).
		GlobalTransitions(gamemachine.Transitions{
			gamemachine.StoppedEvent: gamemachine.StoppedStateIdentifier,
		}).
		FinalStates(gamemachine.StoppedStateIdentifier).
		MustBuild()
}
