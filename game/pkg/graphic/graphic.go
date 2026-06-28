package graphic

import (
	"context"
	"game/pkg/graphic/renderer"
	"game/pkg/graphic/resource"
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
	Shaders() *resource.Shaders
	Telemetry() Telemetry
	Terminate(ctx context.Context)
}
