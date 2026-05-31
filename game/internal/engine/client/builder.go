package client

import (
	"errors"
	"game/pkg/engine"
)

type Builder struct {
	fsm engine.Machine
}

func NewBuilder() Builder {
	return Builder{}
}

func (b Builder) Fsm(machine engine.Machine) Builder {
	b.fsm = machine
	return b
}

func (b Builder) Build() (*Client, error) {
	client := &Client{}

	if b.fsm == nil {
		return nil, errors.New("fsm is nil")
	}
	client.fsm = b.fsm

	return client, nil
}

func (b Builder) MustBuild() *Client {
	cl, err := b.Build()
	if err != nil {
		panic("failed to build client: " + err.Error())
	}

	return cl
}
