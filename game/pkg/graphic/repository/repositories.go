// Package repository provides in-memory stores for graphic resources (shaders
// and uniform blocks).
//
// Repositories are NOT safe for concurrent use. They are intended to be
// accessed from a single goroutine (the render thread); callers are
// responsible for external synchronization if shared across goroutines.
package repository

import (
	gresource "game/pkg/graphic/resource"
	"maps"
)

type UniformBlockRepository struct {
	idIndex   map[uint32]gresource.UniformBlock
	nameIndex map[string]uint32
}

func NewUniformBlockRepository() *UniformBlockRepository {
	return &UniformBlockRepository{
		idIndex:   make(map[uint32]gresource.UniformBlock),
		nameIndex: make(map[string]uint32),
	}
}

func (ub *UniformBlockRepository) ByID(id uint32) (gresource.UniformBlock, bool) {
	uniformBlock, ok := ub.idIndex[id]
	return uniformBlock, ok
}

func (ub *UniformBlockRepository) ByName(name string) (gresource.UniformBlock, bool) {
	id, ok := ub.nameIndex[name]
	if !ok {
		return gresource.UniformBlock{}, false
	}

	return ub.ByID(id)
}

func (ub *UniformBlockRepository) Delete(id uint32) bool {
	uniformBlock, ok := ub.idIndex[id]
	if !ok {
		return false
	}

	delete(ub.idIndex, uniformBlock.ID)
	delete(ub.nameIndex, uniformBlock.Name)

	return true
}

func (ub *UniformBlockRepository) Put(uniformBlock gresource.UniformBlock) bool {
	_, exists := ub.ByID(uniformBlock.ID)
	if exists {
		return false
	}

	ub.idIndex[uniformBlock.ID] = uniformBlock
	ub.nameIndex[uniformBlock.Name] = uniformBlock.ID

	return true
}

func (ub *UniformBlockRepository) All() []gresource.UniformBlock {
	all := make([]gresource.UniformBlock, 0, len(ub.idIndex))
	maps.Values(ub.idIndex)(func(uniformBlock gresource.UniformBlock) bool {
		all = append(all, uniformBlock)
		return true
	})

	return all
}

type ShaderRepository struct {
	idIndex   map[uint32]gresource.Shader
	nameIndex map[string]uint32
}

func NewShaderRepository() *ShaderRepository {
	return &ShaderRepository{
		idIndex:   make(map[uint32]gresource.Shader),
		nameIndex: make(map[string]uint32),
	}
}

func (s *ShaderRepository) ByID(id uint32) (gresource.Shader, bool) {
	shader, ok := s.idIndex[id]
	return shader, ok
}

func (s *ShaderRepository) ByName(name string) (gresource.Shader, bool) {
	id, ok := s.nameIndex[name]
	if !ok {
		return gresource.Shader{}, false
	}

	return s.ByID(id)
}

func (s *ShaderRepository) Delete(id uint32) bool {
	shader, ok := s.idIndex[id]
	if !ok {
		return false
	}

	delete(s.idIndex, shader.ID)
	delete(s.nameIndex, shader.Name)

	return true
}

func (s *ShaderRepository) Put(shader gresource.Shader) bool {
	_, exists := s.ByID(shader.ID)
	if exists {
		return false
	}

	s.idIndex[shader.ID] = shader
	s.nameIndex[shader.Name] = shader.ID

	return true
}

func (s *ShaderRepository) All() []gresource.Shader {
	all := make([]gresource.Shader, 0, len(s.idIndex))
	maps.Values(s.idIndex)(func(shader gresource.Shader) bool {
		all = append(all, shader)
		return true
	})

	return all
}

type Repositories struct {
	Shader  *ShaderRepository
	Uniform *UniformBlockRepository
}
