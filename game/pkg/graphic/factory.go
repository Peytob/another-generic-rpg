package graphic

type Factory interface {
	RenderTarget() RenderTargetFactory
}

type RenderTargetFactory interface {
	WindowRenderTarget() RenderTarget
}
