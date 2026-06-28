package opengl

import (
	"context"
	"fmt"
	"game/pkg/graphic/internal/opengl/resource"
	gresource "game/pkg/graphic/resource"
	"game/pkg/utils/logger"
	"log/slog"

	"github.com/go-gl/gl/v3.3-core/gl"
)

func loadUniformBlocks(ctx context.Context) (*resource.UniformBlocks, error) {
	ub := &resource.UniformBlocks{
		Proj: resource.ProjectionViewUniformBlock{
			UniformBlockDesc: resource.UniformBlockDesc{
				Name:         "ProjView\x00",
				BindingPoint: 0,
			},
			ProjMatrix: resource.UniformBlockVar{
				Name:   "u_proj",
				Size:   16 * 4, // mat4 = 16 floats * 4 bytes
				Offset: 0,
			},
		},
	}

	var projUbo uint32
	gl.GenBuffers(1, &projUbo)
	gl.BindBuffer(gl.UNIFORM_BUFFER, projUbo)
	gl.BufferData(gl.UNIFORM_BUFFER, 16*4, nil, gl.DYNAMIC_DRAW)
	gl.BindBufferBase(gl.UNIFORM_BUFFER, ub.Proj.BindingPoint, projUbo)
	ub.Proj.UniformBlock = resource.Buffer(projUbo)
	logger.FromCtx(ctx).Info("created projection uniform buffer", slog.Int64("id", int64(projUbo)))

	return ub, nil
}

func bindUniformBlock(sp gresource.ShaderProgram, blockName string, bindingPoint uint32) error {
	index := gl.GetUniformBlockIndex(sp.Id(), gl.Str(blockName))
	if index == gl.INVALID_INDEX {
		return fmt.Errorf("uniform block %s index not found", blockName)
	}
	gl.UniformBlockBinding(sp.Id(), index, bindingPoint)
	return nil
}

func bindUniformBlocks(shaders *gresource.Shaders, ub *resource.UniformBlocks) error {
	return bindUniformBlock(shaders.World, ub.Proj.Name, ub.Proj.BindingPoint)
}

func terminateUniformBlocks(ctx context.Context, ub *resource.UniformBlocks) {
	projViewTmp := ub.Proj.UniformBlock.Id()
	logger.FromCtx(ctx).Info("deleting uniform block",
		slog.String("name", ub.Proj.Name),
		slog.Uint64("id", uint64(projViewTmp)))
	gl.DeleteBuffers(1, &projViewTmp)
}
