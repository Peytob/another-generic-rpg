package uniform

import (
	"context"
	"fmt"
	"game/graphic/shader"
	"game/pkg/logger"
	"log/slog"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type Service interface {
	BindBlocksFor(ctx context.Context, program shader.Program) error
	Uniforms() Uniforms
	BlockArray() []Block
	DeleteBlock(ctx context.Context, block Block) error
}

type service struct {
	blockRepository *BlockRepository
	uniforms        *uniforms
}

func InitializeUniforms(ctx context.Context, blockRepository *BlockRepository) (Service, error) {
	s := &service{}

	// uniforms
	s.blockRepository = blockRepository
	s.uniforms = newUniforms(blockRepository)

	if err := s.uniforms.initializeUniformBlocks(ctx); err != nil {
		return nil, fmt.Errorf("initializing uniform blocks: %w", err)
	}

	return s, nil
}

func (s *service) BindBlocksFor(ctx context.Context, program shader.Program) error {
	return s.uniforms.bindBlocksFor(ctx, program)
}

func (s *service) Uniforms() Uniforms {
	return s.uniforms
}

func (s *service) BlockArray() []Block {
	return s.blockRepository.Array()
}

func (s *service) DeleteBlock(ctx context.Context, block Block) error {
	if !gl.IsBuffer(block.ID) {
		return ErrInvalidBlock
	}

	gl.DeleteBuffers(1, &block.ID)

	if !s.blockRepository.Delete(block) {
		logger.FromCtx(ctx).LogAttrs(ctx, slog.LevelWarn,
			"uniform block deleted from ogl but not found in repository", logger.Uint32("uniform_block_id", block.ID))
	}

	logger.FromCtx(ctx).LogAttrs(ctx, slog.LevelInfo, "uniform block deleted", logger.Uint32("uniform_block_id", block.ID))

	return nil
}
