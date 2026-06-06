package graphic

import "context"

// Graphic Describes graphic backend
type Graphic interface {
	NewCanvas() Canvas
	Renderer() Renderer

	Api() Api

	Terminate(ctx context.Context)
}

type Api struct {
	// Graphic backend provider (Vulkan / OGL / etc)
	Name string

	// Graphic backend provider version
	Version string

	// Usually backend hardware name (or other renderer data)
	Renderer string
}
