package rendering

import (
	"context"
	"engine/graphic"
	gbackend "engine/graphic/backend"
	"engine/utils/logger"
	"fmt"
)

func initializeGraphic(ctx context.Context) (graphic.Graphic, error) {
	logger.FromCtx(ctx).Info("initializing graphic")

	graphics, err := gbackend.NewOpenGlGraphics(ctx)
	if err != nil {
		return nil, fmt.Errorf("initialize opengl graphics: %w", err)
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
