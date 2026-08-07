package gamemachine

import (
	"context"
	"engine/graphic"
	"engine/utils/ecs"
	"game/internal/gameplay/tilemap"
	"game/internal/gamestate/ecsentity"
	"game/internal/gamestate/ecssystem"
	"game/internal/rendering/draw"
	"time"
)

const PlayingStateIdentifier = StateIdentifier("playing")

// playingState describes main game state that used while game is running
type playingState struct {
	graphic graphic.Graphic
}

func NewPlayingState(g graphic.Graphic) State {
	return &playingState{
		graphic: g,
	}
}

func (s playingState) Identifier() StateIdentifier {
	return PlayingStateIdentifier
}

func (s playingState) OnEnter(_ context.Context, w ecs.World) error {
	ecsentity.NewTilemapEntity(w, tilemap.MustNewTilemap("123", 3, 64, 64)) // todo mock

	ecsentity.NewRenderStateEntity(w, s.graphic)

	drawer := draw.NewTilemapDrawer(nil)

	w.AddSystem(ecssystem.NewCameraPositionSyncSystem().Execute)
	w.AddSystem(ecssystem.NewTilemapRenderSystem(drawer).Execute)
	w.AddSystem(ecssystem.NewRenderLayersSystem(s.graphic.Repositories().Shader, s.graphic.Renderer()).Execute)

	return nil
}

func (s playingState) OnExit(_ context.Context, _ ecs.World) error {
	return nil
}

func (s playingState) Update(ctx context.Context, world ecs.World, dt time.Duration) (Event, error) {
	if err := world.Update(ctx, dt); err != nil {
		return NoEvent, err
	}
	return NoEvent, nil
}
