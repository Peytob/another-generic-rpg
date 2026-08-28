package ecssystem

import (
	"context"
	"engine/utils/ecs"
	"engine/utils/logger"
	"game/internal/gamestate/ecscomponent"
	"game/internal/rendering/draw"
	"time"
)

type TilemapRender struct {
	tilemapDrawer draw.TilemapDrawer
}

func NewTilemapRender(tilemapDrawer draw.TilemapDrawer) TilemapRender {
	return TilemapRender{
		tilemapDrawer: tilemapDrawer,
	}
}

func (s TilemapRender) Execute(ctx context.Context, world ecs.World, _ time.Duration) error {
	l := logger.FromCtx(ctx).With("system", "tilemap_render")

	tilemapComponent, ok := ecs.GetSingleComponent[ecscomponent.Tilemap](world)
	if !ok {
		l.Warn("no tilemap found, nothing to render")
		return nil
	}

	renderLayersComponent, ok := ecs.GetSingleComponent[ecscomponent.RenderLayers](world)
	if !ok {
		l.Warn("no render layers found")
		return nil
	}

	cameraComponent, ok := ecs.GetSingleComponent[ecscomponent.Camera](world)
	if !ok {
		l.Warn("no camera found")
		return nil
	}

	err := s.tilemapDrawer.Draw(ctx, tilemapComponent.Tilemap, renderLayersComponent.TilemapCanvas, draw.Opts{
		Camera: cameraComponent.Camera,
	})
	if err != nil {
		return err
	}

	return nil
}
