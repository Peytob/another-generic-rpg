package opengl33

import (
	"game/pkg/graphic"
)

type Factory struct {
	renderTarget RenderTargetFactory
}

func NewFactory() Factory {
	return Factory{
		renderTarget: NewRenderTargetFactory(),
	}
}

func (f Factory) RenderTarget() graphic.RenderTargetFactory {
	return f.renderTarget
}

type RenderTargetFactory struct {
}

func NewRenderTargetFactory() RenderTargetFactory {
	return RenderTargetFactory{}
}

func (r RenderTargetFactory) WindowRenderTarget() graphic.RenderTarget {
	return RenderTarget{} // todo
}
