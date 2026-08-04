package client

import (
	"engine/graphic"
	"engine/window"
	"errors"
	"game/internal/gamestate/gamemachine"
)

type Builder struct {
	machine gamemachine.Machine
	window  *window.Window
	graphic graphic.Graphic
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) Machine(machine gamemachine.Machine) *Builder {
	b.machine = machine
	return b
}

func (b *Builder) Window(window *window.Window) *Builder {
	b.window = window
	return b
}

func (b *Builder) Graphic(g graphic.Graphic) *Builder {
	b.graphic = g
	return b
}

func (b *Builder) Build() (*Client, error) {
	if b.machine == nil {
		return nil, errors.New("machine is nil")
	}

	if b.window == nil {
		return nil, errors.New("window is nil")
	}

	if b.graphic == nil {
		return nil, errors.New("graphic is nil")
	}

	return &Client{
		runner:  gamemachine.NewRunner(b.machine),
		window:  b.window,
		graphic: b.graphic,
	}, nil
}

func (b *Builder) MustBuild() *Client {
	cl, err := b.Build()
	if err != nil {
		panic("failed to build client: " + err.Error())
	}

	return cl
}
