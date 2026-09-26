package service

import (
	gresource "engine/graphic/resource"

	"github.com/go-gl/mathgl/mgl32"
)

type UniformBlockAccessor interface {
	GetByID(id uint32) (gresource.UniformBlock, bool)
	GetByName(name string) (gresource.UniformBlock, bool)
	All() []gresource.UniformBlock
}

// UniformBlock service to interact with loaded shaders
type UniformBlock interface {
	UniformBlockAccessor

	CreateUniformBlock(ub gresource.UniformBlock) (gresource.UniformBlock, error)
	SetUniformVariableMat4(ub gresource.UniformBlock, variable string, value mgl32.Mat4) error
	BindBlocksFor(shader gresource.Shader) error
}
