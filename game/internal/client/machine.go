package client

import (
	"game/internal/gamestate/gamemachine"
)

func NewMachine() gamemachine.Machine {
	return gamemachine.NewMachineBuilder().
		InitialState(gamemachine.InitialStateIdentifier).
		BuildState(gamemachine.NewInitialState()).
		Transition(gamemachine.StoppedEvent, gamemachine.StoppedStateIdentifier).
		Build().
		BuildState(gamemachine.NewResourcesLoadingState()).
		Transition(gamemachine.ResourcesLoadedEvent, gamemachine.PlayingStateIdentifier).
		Build().
		BuildState(gamemachine.NewPlayingState()).
		Transition(gamemachine.StoppedEvent, gamemachine.StoppedStateIdentifier).
		Build().
		RegisterState(gamemachine.NewStoppedState(), make(gamemachine.Transitions)).
		GlobalTransitions(gamemachine.Transitions{
			gamemachine.StoppedEvent: gamemachine.StoppedStateIdentifier,
		}).
		FinalStates(gamemachine.StoppedStateIdentifier).
		MustBuild()
}
