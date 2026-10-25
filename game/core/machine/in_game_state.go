package machine

import (
	"context"
	"fmt"
	"game/gameplay/world"
)

const InGameStateIdentifier = "in_game"

type InGameState struct {
	world *world.World
}

func NewInGameState() *InGameState {
	return &InGameState{}
}

func (s *InGameState) Identifier() StateIdentifier {
	return InGameStateIdentifier
}

func (s *InGameState) OnEnter(ctx context.Context) error {
	if err := s.reset(ctx); err != nil {
		return fmt.Errorf("reset: %w", err)
	}

	return nil
}

func (s *InGameState) OnExit(ctx context.Context) error {
	return nil
}

func (s *InGameState) Tick(ctx context.Context, frame Frame) (Event, error) {
	return 0, nil
}

func (s *InGameState) reset(_ context.Context) error {
	s.world = &world.World{}
	return nil
}
