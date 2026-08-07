package client

import (
	"engine/graphic"
	"game/internal/gamestate/gamemachine"
)

func NewMachine(g graphic.Graphic) gamemachine.Machine {
	return gamemachine.NewMachineBuilder().
		InitialState(gamemachine.PlayingStateIdentifier). // todo change to InitialState
		BuildState(gamemachine.NewInitialState()).
		Build().
		BuildState(gamemachine.NewResourcesLoadingState()).
		Transition(gamemachine.ResourcesLoadedEvent, gamemachine.PlayingStateIdentifier).
		Build().
		BuildState(gamemachine.NewPlayingState(g)).
		Build().
		RegisterState(gamemachine.NewStoppedState(), make(gamemachine.Transitions)).
		GlobalTransitions(gamemachine.Transitions{
			gamemachine.StoppedEvent: gamemachine.StoppedStateIdentifier,
		}).
		FinalStates(gamemachine.StoppedStateIdentifier).
		MustBuild()
}
