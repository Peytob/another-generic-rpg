package client

import (
	"game/internal/gamestate/event"
	"game/internal/gamestate/gamemachine"
)

func NewMachine(client *Client) gamemachine.Machine {
	return gamemachine.NewMachineBuilder().
		InitialState(gamemachine.ResourcesLoadingStateIdentifier). // todo change to InitialState
		BuildState(gamemachine.NewInitialState()).
		Build().
		BuildState(gamemachine.NewResourcesLoadingState()).
		Transition(event.ResourcesLoadedEvent, gamemachine.PlayingStateIdentifier).
		Build().
		BuildState(gamemachine.NewPlayingState(client.rendering, client.repositories, client.hid)).
		Build().
		RegisterState(gamemachine.NewStoppedState(), make(gamemachine.Transitions)).
		GlobalTransitions(gamemachine.Transitions{
			event.StoppedEvent: gamemachine.StoppedStateIdentifier,
		}).
		FinalStates(gamemachine.StoppedStateIdentifier).
		MustBuild()
}
