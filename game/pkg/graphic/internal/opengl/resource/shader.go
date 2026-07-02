package resource

import (
	"github.com/go-gl/gl/v3.3-core/gl"
)

type ShaderType uint32

const VertexShader = ShaderType(gl.VERTEX_SHADER)
const GeometryShader = ShaderType(gl.GEOMETRY_SHADER)
const FragmentShader = ShaderType(gl.FRAGMENT_SHADER)

type Shader uint32

func (s Shader) ID() uint32 {
	return uint32(s)
}

func (s Shader) Type() ShaderType {
	var t int32
	gl.GetShaderiv(s.ID(), gl.SHADER_TYPE, &t)
	return ShaderType(t)
}
