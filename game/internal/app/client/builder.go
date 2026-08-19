package client

import (
	"engine/window"
	"errors"
	"game/internal/gamestate/gamemachine"
	"game/internal/gamestate/repositories"
	"game/internal/rendering"
)

type Builder struct {
	machine      gamemachine.Machine
	window       *window.Window
	rendering    rendering.Rendering
	repositories *repositories.Repositories
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

func (b *Builder) Rendering(g rendering.Rendering) *Builder {
	b.rendering = g
	return b
}

func (b *Builder) Repositories(repo repositories.Repositories) *Builder {
	b.repositories = &repo
	return b
}

func (b *Builder) Build() (*Client, error) {
	if b.machine == nil {
		return nil, errors.New("machine is nil")
	}

	if b.window == nil {
		return nil, errors.New("window is nil")
	}

	if b.rendering == nil {
		return nil, errors.New("rendering is nil")
	}

	if b.repositories == nil {
		return nil, errors.New("repositories is nil")
	}

	return &Client{
		runner:       gamemachine.NewRunner(b.machine),
		window:       b.window,
		rendering:    b.rendering,
		repositories: *b.repositories,
	}, nil
}

func (b *Builder) MustBuild() *Client {
	cl, err := b.Build()
	if err != nil {
		panic("failed to build client: " + err.Error())
	}

	return cl
}
