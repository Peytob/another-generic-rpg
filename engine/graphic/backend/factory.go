package backend

import (
	"context"
	"engine/graphic"
	"engine/graphic/internal/opengl"
)

func NewOpenGlGraphics(ctx context.Context) (graphic.Graphic, error) {
	return opengl.NewGraphic(ctx)
}
