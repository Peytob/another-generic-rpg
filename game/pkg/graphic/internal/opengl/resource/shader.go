package resource

import (
	"context"
	"fmt"

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

func BuildShader(_ context.Context, code string, shaderType ShaderType) (Shader, error) {
	codeStrPtr, freeCodeStrings := gl.Strs(code)
	defer freeCodeStrings()

	shader := gl.CreateShader(uint32(shaderType))
	gl.ShaderSource(shader, 1, codeStrPtr, nil)
	gl.CompileShader(shader)

	var success int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &success)
	if success == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := make([]uint8, logLength)
		var actualLength int32
		gl.GetShaderInfoLog(shader, logLength, &actualLength, &log[0])

		gl.DeleteShader(shader)
		return 0, fmt.Errorf("failed to compile shader: %s", string(log))
	}

	return Shader(shader), nil
}
