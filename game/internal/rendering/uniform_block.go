package rendering

import (
	"context"
	"engine/graphic"
	gresource "engine/graphic/resource"
	"engine/utils/logger"
)

const (
	ProjViewUniformBlock = "ProjView"
	ProjUniform          = "u_proj"
)

func loadUniformBlocks(ctx context.Context, g graphic.Graphic) error {
	log := logger.FromCtx(ctx)
	log.Info("creating uniform blocks")

	projViewBlock := gresource.UniformBlock{
		Name:         ProjViewUniformBlock,
		BindingPoint: 0,

		Variables: gresource.NewUniformVariables().
			Set(ProjUniform, 16*4, 0), // mat4 = 16 floats * 4 bytes,
	}

	projViewBlock, err := g.Services().Uniform.CreateUniformBlock(projViewBlock)
	if err != nil {
		return err
	}
	log.Info("created uniform block", "name", projViewBlock.Name, "id", projViewBlock.ID)

	return nil
}

func bindUniformBlocks(ctx context.Context, g graphic.Graphic) error {
	logger.FromCtx(ctx).Info("binding uniform blocks")

	for _, shader := range g.Repositories().Shader.All() {
		err := g.Services().Uniform.BindBlocksFor(shader)
		if err != nil {
			return err
		}
	}

	return nil
}
