package graphic

import "context"

// Graphic describes graphic backend. Graphic will provide all methods to interaction with graphical backend, includes
// factory methods for shaders, view points and other graphical stuff
type Graphic interface {
	// NewCanvas creates new empty canvas. Canvas will be compatible with graphical backend
	NewCanvas() Canvas

	// Renderer returns configured and ready renderer
	Renderer() Renderer

	// Factory used for creating backend-compatible graphic resources
	Factory() Factory

	// Telemetry returns information about graphical backend
	Telemetry() Telemetry

	// Terminate terminates graphics instance
	Terminate(ctx context.Context)
}

type Telemetry struct {
	// Graphic backend provider (Vulkan / OGL / etc)
	Name string

	// Graphic backend provider version
	Version string

	// Usually backend hardware name (or other renderer data)
	Renderer string
}
