package rendering

import (
	"context"
	"engine/graphic"
	grenderer "engine/graphic/renderer"
	grepository "engine/graphic/repository"
	gservice "engine/graphic/service"
	"game/internal/gamestate/repositories"
)

// Rendering extension for graphic module with game-specific logic
type Rendering interface {
	graphic.Graphic

	Drawers() *Drawers
}

type rendering struct {
	graphic graphic.Graphic

	drawers *Drawers
}

func NewRendering(ctx context.Context, repo repositories.Repositories) (Rendering, error) {
	g, err := initializeGraphic(ctx)
	if err != nil {
		return nil, err
	}

	return &rendering{
		graphic: g,
		drawers: NewDrawers(g, repo),
	}, nil
}

func (r rendering) Renderer() grenderer.Renderer {
	return r.graphic.Renderer()
}

func (r rendering) NewCanvas() grenderer.Canvas {
	return r.graphic.NewCanvas()
}

func (r rendering) Services() gservice.Services {
	return r.graphic.Services()
}

func (r rendering) Repositories() grepository.Repositories {
	return r.graphic.Repositories()
}

func (r rendering) Telemetry() graphic.Telemetry {
	return r.graphic.Telemetry()
}

func (r rendering) Terminate(ctx context.Context) {
	r.graphic.Terminate(ctx)
}

func (r rendering) Drawers() *Drawers {
	return r.drawers
}
