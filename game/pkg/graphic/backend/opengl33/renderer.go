package opengl33

import (
	"context"
	"errors"
	"game/pkg/graphic"
)

type RenderTarget struct {
}

type Renderer struct {
	defaultRenderTarget RenderTarget
}

func NewRenderer(defaultRenderTarget RenderTarget) Renderer {
	return Renderer{
		defaultRenderTarget: defaultRenderTarget,
	}
}

func (r Renderer) Render(ctx context.Context, canvas graphic.Canvas, opts graphic.RenderOpts) error {
	_, ok := canvas.(Canvas)
	if !ok {
		return errors.New("canvas is not OGL-compilable")
	}

	if opts.RenderTarget == nil {
		opts.RenderTarget = r.defaultRenderTarget
	}

	// TODO
	//if opts.Shader == nil {
	//	return errors.New("shader is not defined")
	//}

	_, ok = opts.RenderTarget.(RenderTarget)
	if !ok {
		return errors.New("renderTarget OGL-compilable")
	}

	return nil
}
