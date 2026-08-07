package ecssystem

import (
	"context"
	"engine/utils/ecs"
	"engine/utils/logger"
	"game/internal/gamestate/ecscomponent"
	"game/internal/rendering/draw"
	"time"
)

type TilemapRenderSystem struct {
	tilemapDrawer draw.TilemapDrawer
}

func NewTilemapRenderSystem(tilemapDrawer draw.TilemapDrawer) TilemapRenderSystem {
	return TilemapRenderSystem{
		tilemapDrawer: tilemapDrawer,
	}
}

func (s TilemapRenderSystem) Execute(ctx context.Context, world ecs.World, _ time.Duration) error {
	l := logger.FromCtx(ctx).With("system", "tilemap_render")

	tilemapComponent, ok := ecs.GetSingleComponent[ecscomponent.TilemapComponent](world)
	if !ok {
		l.Warn("no tilemap found, nothing to render")
		return nil
	}

	renderLayersComponent, ok := ecs.GetSingleComponent[ecscomponent.RenderLayersComponent](world)
	if !ok {
		l.Warn("no render layers found")
		return nil
	}

	cameraComponent, ok := ecs.GetSingleComponent[ecscomponent.CameraComponent](world)
	if !ok {
		l.Warn("no camera found")
		return nil
	}

	err := s.tilemapDrawer.Draw(ctx, tilemapComponent.Tilemap, renderLayersComponent.TilemapCanvas, draw.DrawOpts{
		Camera: cameraComponent.Camera,
	})
	if err != nil {
		return err
	}

	return nil
}
