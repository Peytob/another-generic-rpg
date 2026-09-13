package state

import (
	"game/internal/game/repositories"
	"game/internal/gameplay/world"
)

type State struct {
	Repositories repositories.Repositories
	World        *world.World
}
