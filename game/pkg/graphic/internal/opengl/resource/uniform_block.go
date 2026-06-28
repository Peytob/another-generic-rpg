package resource

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type UniformBlocks struct {
	Proj ProjectionViewUniformBlock
}

type UniformBlockVar struct {
	Name   string
	Size   int
	Offset int
}

type UniformBlockDesc struct {
	UniformBlock Buffer

	Name         string
	BindingPoint uint32
}

type ProjectionViewUniformBlock struct {
	UniformBlockDesc

	ProjMatrix UniformBlockVar
}

func (ub *ProjectionViewUniformBlock) SetProjectionMatrix(m mgl32.Mat4) {
	uniformBlockMat4(ub.UniformBlock, ub.ProjMatrix.Offset, m)
}

func uniformBlockMat4(ub Buffer, offset int, mat mgl32.Mat4) {
	gl.BindBuffer(gl.UNIFORM_BUFFER, ub.Id())
	gl.BufferSubData(gl.UNIFORM_BUFFER, offset, 16*4, gl.Ptr(&mat[0]))
}
