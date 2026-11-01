package uniform

import (
	"context"
	"errors"
	"fmt"
	"game/graphic/shader"
	"game/pkg/logger"
	"log/slog"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	projViewUniformBlock = "ProjView"
	projUniform          = "u_proj"
)

type Uniforms interface {
	ProjectionMatrix(proj mgl32.Mat4) error
}

type uniforms struct {
	blockRepository *BlockRepository
}

func newUniforms(blockRepository *BlockRepository) *uniforms {
	return &uniforms{
		blockRepository: blockRepository,
	}
}

func (u *uniforms) ProjectionMatrix(proj mgl32.Mat4) error {
	block, ok := u.blockRepository.GetByName(projViewUniformBlock)
	if !ok {
		return fmt.Errorf("block '%s' not found in repository", projViewUniformBlock)
	}

	return setUniformVariableMat4(block, projUniform, proj)
}

func (u *uniforms) createUniformBlock(ctx context.Context, uniformBlock Block) (Block, error) {
	_, exists := u.blockRepository.GetByName(uniformBlock.Name)
	if exists {
		return Block{}, fmt.Errorf("uniform block with given name already registered")
	}

	totalSize := uniformBlock.Variables.TotalSize()

	var ubo uint32
	gl.GenBuffers(1, &ubo)
	gl.BindBuffer(gl.UNIFORM_BUFFER, ubo)
	gl.BufferData(gl.UNIFORM_BUFFER, totalSize, nil, gl.DYNAMIC_DRAW)
	gl.BindBufferBase(gl.UNIFORM_BUFFER, uniformBlock.BindingPoint, ubo)

	uniformBlock.ID = ubo

	if !u.blockRepository.Put(uniformBlock) {
		gl.DeleteBuffers(1, &ubo)
		return Block{}, fmt.Errorf("uniform block already registered")
	}

	logger.FromCtx(ctx).LogAttrs(ctx, slog.LevelInfo, "created uniform block",
		logger.Uint32("uniform_block_id", uniformBlock.ID),
		slog.String("uniform_block_name", uniformBlock.Name))

	return uniformBlock, nil
}

func (u *uniforms) initializeUniformBlocks(ctx context.Context) error {
	logger.FromCtx(ctx).LogAttrs(ctx, slog.LevelInfo, "initializing uniform blocks")
	var creatingErrors, err error

	_, err = u.createUniformBlock(ctx, Block{
		Name:         projViewUniformBlock,
		BindingPoint: 0,
		Variables: newUniformVariables().
			Set(projUniform, 16*4, 0), // mat4 = 16 floats * 4 bytes,
	})
	creatingErrors = errors.Join(creatingErrors, err)

	return creatingErrors
}

func (u *uniforms) bindBlocksFor(ctx context.Context, program shader.Program) error {
	var err error

	var count int32
	gl.GetProgramiv(program.ID, gl.ACTIVE_UNIFORM_BLOCKS, &count)

	for i := uint32(0); i < uint32(count); i++ {
		var nameLen int32
		var nameBuf [256]byte
		gl.GetActiveUniformBlockName(program.ID, i, int32(len(nameBuf)), &nameLen, &nameBuf[0])
		name := string(nameBuf[:nameLen])

		ub, ok := u.blockRepository.GetByName(name)
		if !ok {
			err = errors.Join(err, fmt.Errorf("uniform block %s from shader program %s not found in registered uniform blocks",
				name, program.Name))
			continue
		}

		gl.UniformBlockBinding(program.ID, i, ub.BindingPoint)

		logger.FromCtx(ctx).LogAttrs(ctx, slog.LevelInfo, "bind uniform block for shader program",
			logger.Uint32("uniform_block_id", ub.ID),
			logger.Uint32("program_id", program.ID),
			logger.Uint32("binding_point", ub.BindingPoint))
	}

	return err
}
