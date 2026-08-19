package gamemachine

import (
	"context"
	"engine/utils/ecs"
	"game/internal/gameplay/tilemap"
	"game/internal/gamestate/ecsentity"
	"game/internal/gamestate/ecssystem"
	"game/internal/gamestate/repositories"
	"game/internal/rendering"
	"time"
)

const PlayingStateIdentifier = StateIdentifier("playing")

// playingState describes main game state that used while game is running
type playingState struct {
	renderer       rendering.Rendering
	tileRepository *repositories.TileRepository
}

func NewPlayingState(r rendering.Rendering, repo repositories.Repositories) State {
	return &playingState{
		renderer:       r,
		tileRepository: repo.TileRepository,
	}
}

func (s playingState) Identifier() StateIdentifier {
	return PlayingStateIdentifier
}

func (s playingState) OnEnter(_ context.Context, w ecs.World) error {
	ecsentity.NewTilemapEntity(w, tilemap.MustNewTilemap("123", 3, 64, 64)) // todo mock

	ecsentity.NewRenderStateEntity(w, s.renderer)

	drawer := s.renderer.Drawers().Tilemap

	w.AddSystem(ecssystem.NewCameraPositionSyncSystem())
	w.AddSystem(ecssystem.NewTilemapRenderSystem(drawer))
	w.AddSystem(ecssystem.NewRenderLayersSystem(s.renderer.Repositories().Shader, s.renderer.Renderer()))

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
