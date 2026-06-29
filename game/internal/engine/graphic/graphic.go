package graphic

import (
	"context"
	"game/pkg/graphic"
	gbackend "game/pkg/graphic/backend"
	"game/pkg/utils/logger"
)

func InitializeGraphic(ctx context.Context) (graphic.Graphic, error) {
	logger.FromCtx(ctx).Info("initializing graphic")

	graphics, err := gbackend.NewOpenGlGraphics(ctx)
	if err != nil {
		return nil, err
	}

	err = loadShaders(ctx, graphics)
	if err != nil {
		return nil, err
	}

	err = loadUniformBlocks(ctx, graphics)
	if err != nil {
		return nil, err
	}

	err = bindUniformBlocks(ctx, graphics)
	if err != nil {
		return nil, err
	}

	return graphics, nil
}
