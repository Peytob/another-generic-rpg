package backend

import (
	"context"
	"game/pkg/graphic"
	"game/pkg/graphic/internal/opengl"
)

func NewOpenGlGraphics(ctx context.Context) (graphic.Graphic, error) {
	return opengl.NewGraphic(ctx)
}
