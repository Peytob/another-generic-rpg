package client

import (
	"game/internal/game/machine"
	"game/internal/game/state"
)

func NewMachine(client *Client) machine.Machine {
	return machine.NewBuilder().
		InitialState(machine.ResourcesLoadingStateIdentifier). // todo change to InitialState
		BuildState(machine.NewInitialState()).
		Build().
		BuildState(machine.NewResourcesLoadingState()).
		Transition(state.ResourcesLoadedEvent, machine.PlayingStateIdentifier).
		Build().
		BuildState(machine.NewPlayingState(client.rendering, client.repositories, client.hid)).
		Build().
		RegisterState(machine.NewStoppedState(), make(machine.Transitions)).
		GlobalTransitions(machine.Transitions{
			state.StoppedEvent: machine.StoppedStateIdentifier,
		}).
		FinalStates(machine.StoppedStateIdentifier).
		MustBuild()
}
