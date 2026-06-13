package resource

import (
	"context"
	"fmt"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type ShaderProgram uint32

func (sp ShaderProgram) Id() uint32 {
	return uint32(sp)
}

func (sp ShaderProgram) Delete() {
	gl.DeleteShader(sp.Id())
}

type ShaderProgramBuilder struct {
	Vertex   Shader
	Fragment Shader
}

func BuildShaderProgram(_ context.Context, builder ShaderProgramBuilder) (ShaderProgram, error) {
	if builder.Vertex.Id() == 0 {
		return 0, fmt.Errorf("vertex shader required to create shader program")
	}

	if builder.Vertex.Type() != VertexShader {
		return 0, fmt.Errorf("vertex shader has wrong shader type")
	}

	if builder.Fragment.Id() == 0 {
		return 0, fmt.Errorf("fragment shader required to create shader program")
	}

	if builder.Fragment.Type() != FragmentShader {
		return 0, fmt.Errorf("fragment shader has wrong shader type")
	}

	shaderProgram := gl.CreateProgram()
	gl.AttachShader(shaderProgram, builder.Vertex.Id())
	gl.AttachShader(shaderProgram, builder.Fragment.Id())
	gl.LinkProgram(shaderProgram)

	var success int32
	gl.GetProgramiv(shaderProgram, gl.LINK_STATUS, &success)
	if success == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(shaderProgram, gl.INFO_LOG_LENGTH, &logLength)

		log := make([]uint8, logLength)
		var actualLength int32
		gl.GetProgramInfoLog(shaderProgram, logLength, &actualLength, &log[0])

		gl.DeleteProgram(shaderProgram)

		return 0, fmt.Errorf("failed to link shader: %s", string(log))
	}

	return ShaderProgram(shaderProgram), nil
}
