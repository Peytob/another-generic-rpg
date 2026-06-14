package resource

import (
	"context"
	"fmt"
	"game/pkg/math"
	"strings"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type ShaderProgram uint32

func (sp ShaderProgram) Id() uint32 {
	return uint32(sp)
}

func (sp ShaderProgram) Delete() {
	gl.DeleteShader(sp.Id())
}

func (sp ShaderProgram) UniformMat4(variable string, mat mgl32.Mat4) error {
	location := sp.getUniformLocation(variable)
	if location == -1 {
		return fmt.Errorf("uniform %s not found", variable)
	}

	gl.UniformMatrix4fv(location, 1, false, &mat[0])
	return nil
}

func (sp ShaderProgram) UniformTransform(variable string, transform math.Transform) error {
	return sp.UniformMat4(variable, mgl32.Mat4(transform))
}

func (sp ShaderProgram) getUniformLocation(variable string) int32 {
	if !strings.HasSuffix(variable, "\x00") {
		variable = variable + "\x00"
	}
	return gl.GetUniformLocation(sp.Id(), gl.Str(variable))
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
