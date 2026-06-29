package service

import (
	"context"
	gresource "game/pkg/graphic/resource"
	"game/pkg/math"

	"github.com/go-gl/mathgl/mgl32"
)

// ShaderLoader loader for GLSL shaders and programs / pipelines. Loader should automatically save created shader
// into repository and backend.
type ShaderLoader interface {
	// LoadGlslStageShader compiles GLSL shader.
	LoadGlslStageShader(code string, stage gresource.ShaderStage) (gresource.LoadedShaderStage, error)

	// BuildShaderProgram builds shader program from builder
	BuildShaderProgram(shaderBuilder gresource.ShaderBuilder, name string) (gresource.Shader, error)

	// Terminate cleanup state after loading shaders
	Terminate(ctx context.Context)
}

// Shader service to interact with loaded shaders
type Shader interface {
	NewShaderLoader() ShaderLoader
	UniformTransform(sp gresource.Shader, variable string, transform math.Transform) error
	UniformMat4(sp gresource.Shader, variable string, mat mgl32.Mat4) error
}
