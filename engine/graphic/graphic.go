package graphic

import (
	"context"
	"engine/graphic/renderer"
	"engine/graphic/repository"
	"engine/graphic/service"
)

type Telemetry struct {
	// Graphic backend provider (Vulkan / OGL / etc)
	Name string

	// Graphic backend provider version
	Version string

	// Usually backend hardware name (or other renderer data)
	Renderer string
}

type Graphic interface {
	Renderer() renderer.Renderer
	NewCanvas() renderer.Canvas

	Services() service.Services
	Repositories() repository.Repositories

	Telemetry() Telemetry
	Terminate(ctx context.Context)
}
