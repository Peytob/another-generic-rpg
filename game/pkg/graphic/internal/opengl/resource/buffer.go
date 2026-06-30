package resource

import "github.com/go-gl/gl/v3.3-core/gl"

type VertexArray uint32

func (vao VertexArray) ID() uint32 {
	return uint32(vao)
}

func (vao VertexArray) Delete() {
	id := vao.ID()
	gl.DeleteVertexArrays(1, &id)
}

type Buffer uint32

func (b Buffer) ID() uint32 {
	return uint32(b)
}
