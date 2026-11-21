package machine

import (
	"context"
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

func (s *InGameState) OnEnter(_ context.Context) error {
	s.reset()

	return nil
}

func (s *InGameState) OnExit(_ context.Context) error {
	return nil
}

func (s *InGameState) Tick(_ context.Context, _ Frame) (Event, error) {
	return 0, nil
}

func (s *InGameState) reset() {
	s.world = &world.World{}
}
