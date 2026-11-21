package client

import (
	"game/core/machine"
	"game/graphic"
	"game/window"
)

type Builder struct {
	window  *window.Window
	graphic *graphic.Graphic
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) Window(window *window.Window) *Builder {
	b.window = window
	return b
}

func (b *Builder) Graphic(g *graphic.Graphic) *Builder {
	b.graphic = g
	return b
}

func (b *Builder) Build() (*Client, error) {
	client := &Client{}

	if b.window == nil {
		return nil, ErrWindowIsNil
	}
	client.window = b.window

	if b.graphic == nil {
		return nil, ErrGraphicIsNil
	}
	client.graphic = b.graphic

	initState := machine.NewInitState()
	stoppedState := machine.NewStoppedState()

	inGameState := machine.NewInGameState()
	resourcesLoadingState := machine.NewResourcesLoadingState()

	client.machine = machine.NewBuilder().
		BuildState(initState).
		Transition(machine.InitializedEvent, machine.ResourcesLoadingStateIdentifier).
		Build().
		BuildState(resourcesLoadingState).
		Transition(machine.ResourcesLoadedEvent, machine.InGameStateIdentifier).
		Build().
		BuildState(stoppedState).Build().
		BuildState(inGameState).Build().
		InitialState(initState.Identifier()).
		FinalStates(stoppedState.Identifier()).
		GlobalTransitions(machine.Transitions{
			machine.StoppedEvent: stoppedState.Identifier(),
		}).
		MustBuild()

	return client, nil
}

func (b *Builder) MustBuild() *Client {
	cl, err := b.Build()
	if err != nil {
		panic("failed to build client: " + err.Error())
	}

	return cl
}
