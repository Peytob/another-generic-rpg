package service

import (
	"game/pkg/graphic/resource"

	"github.com/go-gl/mathgl/mgl32"
)

// UniformBlock service to interact with loaded shaders
type UniformBlock interface {
	CreateUniformBlock(ub *resource.UniformBlock) error
	SetUniformVariableMat4(ub resource.UniformBlock, variable string, value mgl32.Mat4) error
	BindBlocksFor(shader resource.Shader) error
}
