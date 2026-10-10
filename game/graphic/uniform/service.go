package uniform

import (
	"context"
	"fmt"
	"game/graphic/shader"
)

type Service interface {
	BindBlocksFor(ctx context.Context, program shader.Program) error
	Uniforms() Uniforms
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
