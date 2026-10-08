package uniform

import (
	"context"
	"errors"
	"fmt"
	"game/graphic/shader"
	"game/utils/logger"
	"log/slog"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Service interface {
	CreateUniformBlock(ctx context.Context, uniformBlock Block) (Block, error)
	SetUniformVariableMat4(ub Block, variableName string, mat mgl32.Mat4) error
	BindBlocksFor(ctx context.Context, program shader.Program) error
	GetBlockByID(id uint32) (Block, bool)
	GetBlockByName(name string) (Block, bool)
	BlockArray() []Block
}

type service struct {
	blockRepository *BlockRepository
}

func NewUniformService(blockRepository *BlockRepository) Service {
	return &service{
		blockRepository: blockRepository,
	}
}

func (s *service) CreateUniformBlock(ctx context.Context, uniformBlock Block) (Block, error) {
	_, exists := s.blockRepository.GetByName(uniformBlock.Name)
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

	if !s.blockRepository.Put(uniformBlock) {
		gl.DeleteBuffers(1, &ubo)
		return Block{}, fmt.Errorf("uniform block already registered")
	}

	logger.FromCtx(ctx).LogAttrs(ctx, slog.LevelInfo, "created uniform block",
		logger.Uint32("uniform_block_id", uniformBlock.ID),
		slog.String("uniform_block_name", uniformBlock.Name))

	return uniformBlock, nil
}

func (s *service) SetUniformVariableMat4(uniformBlock Block, variableName string, mat mgl32.Mat4) error {
	variable, found := uniformBlock.Variables.Lookup(variableName)
	if !found {
		return fmt.Errorf("variable not found")
	}

	gl.BindBuffer(gl.UNIFORM_BUFFER, uniformBlock.ID)
	gl.BufferSubData(gl.UNIFORM_BUFFER, variable.Offset, variable.Size, gl.Ptr(&mat[0]))

	return nil
}

func (s *service) BindBlocksFor(ctx context.Context, program shader.Program) error {
	var err error

	var count int32
	gl.GetProgramiv(program.ID, gl.ACTIVE_UNIFORM_BLOCKS, &count)

	for i := uint32(0); i < uint32(count); i++ {
		var nameLen int32
		var nameBuf [256]byte
		gl.GetActiveUniformBlockName(program.ID, i, int32(len(nameBuf)), &nameLen, &nameBuf[0])
		name := string(nameBuf[:nameLen])

		ub, ok := s.blockRepository.GetByName(name)
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

func (s *service) GetBlockByID(id uint32) (Block, bool) {
	return s.blockRepository.GetByID(id)
}

func (s *service) GetBlockByName(name string) (Block, bool) {
	return s.blockRepository.GetByName(name)
}

func (s *service) BlockArray() []Block {
	return s.blockRepository.Array()
}
