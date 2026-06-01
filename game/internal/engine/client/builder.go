package client

import (
	"errors"
	"game/pkg/engine"
	"game/pkg/window"
)

type Builder struct {
	fsm    engine.Machine
	window *window.Window
}

func NewBuilder() Builder {
	return Builder{}
}

func (b Builder) Fsm(machine engine.Machine) Builder {
	b.fsm = machine
	return b
}

func (b Builder) Window(window *window.Window) Builder {
	b.window = window
	return b
}

func (b Builder) Build() (*Client, error) {
	client := &Client{}

	if b.fsm == nil {
		return nil, errors.New("fsm is nil")
	}
	client.fsm = b.fsm

	if b.window == nil {
		return nil, errors.New("fsm is nil")
	}
	client.window = b.window

	return client, nil
}

func (b Builder) MustBuild() *Client {
	cl, err := b.Build()
	if err != nil {
		panic("failed to build client: " + err.Error())
	}

	return cl
}
