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

type RenderLayersSystem struct {
	shaderRepository *repository.ShaderRepository
	renderer         renderer.Renderer
}

func NewRenderLayersSystem(shaderRepository *repository.ShaderRepository, renderer renderer.Renderer) RenderLayersSystem {
	return RenderLayersSystem{
		shaderRepository: shaderRepository,
		renderer:         renderer,
	}
}

func (s RenderLayersSystem) Execute(ctx context.Context, world ecs.World, _ time.Duration) error {
	l := logger.FromCtx(ctx).With("system", "render_layers")

	renderLayersComponent, ok := ecs.GetSingleComponent[ecscomponent.RenderLayersComponent](world)
	if !ok {
		l.Warn("no render layers found")
		return nil
	}

	shader, shaderFound := s.shaderRepository.ByName(rendering.TilemapShader)
	if !shaderFound {
		return fmt.Errorf("no shader found in graphic repositories")
	}

	cameraComponent, ok := ecs.GetSingleComponent[ecscomponent.CameraComponent](world)
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
