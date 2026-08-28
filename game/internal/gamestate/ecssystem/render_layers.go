package ecssystem

import (
	"context"
	"engine/graphic/renderer"
	"engine/graphic/repository"
	"engine/math"
	"engine/utils/ecs"
	"engine/utils/logger"
	"fmt"
	"game/internal/gamestate/ecscomponent"
	"game/internal/rendering"
	"time"
)

type RenderLayers struct {
	shaderRepository *repository.ShaderRepository
	renderer         renderer.Renderer
}

func NewRenderLayers(shaderRepository *repository.ShaderRepository, renderer renderer.Renderer) RenderLayers {
	return RenderLayers{
		shaderRepository: shaderRepository,
		renderer:         renderer,
	}
}

func (s RenderLayers) Execute(ctx context.Context, world ecs.World, _ time.Duration) error {
	l := logger.FromCtx(ctx).With("system", "render_layers")

	renderLayersComponent, ok := ecs.GetSingleComponent[ecscomponent.RenderLayers](world)
	if !ok {
		l.Warn("no render layers found")
		return nil
	}

	shader, shaderFound := s.shaderRepository.ByName(rendering.TilemapShader)
	if !shaderFound {
		return fmt.Errorf("no shader found in graphic repositories")
	}

	cameraComponent, ok := ecs.GetSingleComponent[ecscomponent.Camera](world)
	if !ok {
		l.Warn("no camera found")
		return nil
	}

	opts := renderer.RenderOpts{
		View: cameraComponent.Camera.ViewMatrix(),
		Model: math.NewTransformation().
			Transform(),
		Shader: shader,
		Mode:   renderer.Wireframe,
	}

	err := s.renderer.Render(ctx, renderLayersComponent.TilemapCanvas, opts)
	if err != nil {
		return fmt.Errorf("failed to render: %w", err)
	}

	return nil
}
