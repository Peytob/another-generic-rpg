package service

import (
	"fmt"
	oglresource "game/pkg/graphic/internal/opengl/resource"
	grepository "game/pkg/graphic/repository"
	gresource "game/pkg/graphic/resource"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type UniformBlock struct {
	uniformBlocks *grepository.UniformBlockRepository
}

func NewUniformBlock(ub *grepository.UniformBlockRepository) *UniformBlock {
	return &UniformBlock{
		uniformBlocks: ub,
	}
}

func (u *UniformBlock) CreateUniformBlock(ub *gresource.UniformBlock) error {
	totalSize := 0
	for _, variable := range ub.Variables.All() {
		totalSize += variable.Size
	}

	var projUbo uint32
	gl.GenBuffers(1, &projUbo)
	gl.BindBuffer(gl.UNIFORM_BUFFER, projUbo)
	gl.BufferData(gl.UNIFORM_BUFFER, totalSize, nil, gl.DYNAMIC_DRAW)
	gl.BindBufferBase(gl.UNIFORM_BUFFER, ub.BindingPoint, projUbo)

	ub.ID = projUbo

	return nil
}

func (u *UniformBlock) SetUniformVariableMat4(ub gresource.UniformBlock, variableName string, mat mgl32.Mat4) error {
	variable, found := ub.Variables.Lookup(variableName)
	if !found {
		return fmt.Errorf("variable not found")
	}

	gl.BindBuffer(gl.UNIFORM_BUFFER, ub.ID)
	gl.BufferSubData(gl.UNIFORM_BUFFER, variable.Offset, variable.Size, gl.Ptr(&mat[0]))
	return nil
}

func (u *UniformBlock) BindBlocksFor(shader gresource.Shader) error {
	shaderProgram := oglresource.ShaderProgram(shader.ID)

	var count int32
	gl.GetProgramiv(shaderProgram.ID(), gl.ACTIVE_UNIFORM_BLOCKS, &count)

	for i := uint32(0); i < uint32(count); i++ {
		var nameLen int32
		var nameBuf [256]byte
		gl.GetActiveUniformBlockName(shaderProgram.ID(), i, int32(len(nameBuf)), &nameLen, &nameBuf[0])
		name := string(nameBuf[:nameLen])

		ub, ok := u.uniformBlocks.ByName(name)
		if !ok {
			return fmt.Errorf("uniform block %s from shader program %s not found in registered uniform blocks", name, shader.Name)
		}

		gl.UniformBlockBinding(shaderProgram.ID(), i, ub.BindingPoint)
	}

	return nil
}
