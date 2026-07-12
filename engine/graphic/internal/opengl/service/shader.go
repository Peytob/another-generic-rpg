package service

import (
	"context"
	oglresource "engine/graphic/internal/opengl/resource"
	grepository "engine/graphic/repository"
	gresource "engine/graphic/resource"
	gservice "engine/graphic/service"
	"engine/math"
	"engine/utils/logger"
	"fmt"
	"log/slog"
	"strings"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type shaderLoader struct {
	shaders        map[uint32]oglresource.Shader
	shaderPrograms *grepository.ShaderRepository
}

func newShaderLoader(sr *grepository.ShaderRepository) *shaderLoader {
	return &shaderLoader{
		shaders:        make(map[uint32]oglresource.Shader),
		shaderPrograms: sr,
	}
}

func (s *shaderLoader) LoadGlslStageShader(code string, stage gresource.ShaderStage) (gresource.LoadedShaderStage, error) {
	shaderType := shaderStageToShaderType(stage)
	if shaderType == 0 {
		return gresource.LoadedShaderStage{}, fmt.Errorf("unknown or unsupported shader stage")
	}

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
		return gresource.LoadedShaderStage{}, fmt.Errorf("failed to compile shader: %s", string(log))
	}

	s.shaders[shader] = oglresource.Shader(shader)

	return gresource.LoadedShaderStage{
		ID:    shader,
		Stage: stage,
	}, nil
}

func (s *shaderLoader) BuildShaderProgram(builder *gresource.ShaderBuilder, name string) (gresource.Shader, error) {
	vertexID := builder.Stage(gresource.VertexStage).ID

	if vertexID == 0 {
		return gresource.Shader{}, fmt.Errorf("vertex shader required to create shader program")
	}

	vertexShader, ok := s.shaders[vertexID]
	if !ok {
		return gresource.Shader{}, fmt.Errorf("unknown vertex shader")
	}

	if vertexShader.Type() != oglresource.VertexShader {
		return gresource.Shader{}, fmt.Errorf("vertex shader has wrong shader type")
	}

	fragmentID := builder.Stage(gresource.FragmentStage).ID

	if fragmentID == 0 {
		return gresource.Shader{}, fmt.Errorf("fragment shader required to create shader program")
	}

	fragmentShader, ok := s.shaders[fragmentID]
	if !ok {
		return gresource.Shader{}, fmt.Errorf("unknown fragment shader")
	}

	if fragmentShader.Type() != oglresource.FragmentShader {
		return gresource.Shader{}, fmt.Errorf("fragment shader has wrong shader type")
	}

	shaderProgram := gl.CreateProgram()
	gl.AttachShader(shaderProgram, vertexShader.ID())
	gl.AttachShader(shaderProgram, fragmentShader.ID())
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

		return gresource.Shader{}, fmt.Errorf("failed to link shader: %s", string(log))
	}

	program := gresource.Shader{
		ID:   shaderProgram,
		Name: name,
	}

	if !s.shaderPrograms.Put(program) {
		gl.DeleteProgram(shaderProgram)
		return gresource.Shader{}, fmt.Errorf("shader program %s already registered", name)
	}

	return program, nil
}

func (s *shaderLoader) Terminate(ctx context.Context) {
	for _, v := range s.shaders {
		logger.FromCtx(ctx).Info("deleting shader", slog.Uint64("id", uint64(v.ID())))
		gl.DeleteShader(v.ID())
	}
}

type Shader struct {
	shaderRepository *grepository.ShaderRepository
}

func NewShader(sr *grepository.ShaderRepository) *Shader {
	return &Shader{
		shaderRepository: sr,
	}
}

func (s *Shader) UniformMat4(sp gresource.Shader, variable string, mat mgl32.Mat4) error {
	location := getUniformLocation(sp.ID, variable)
	if location == -1 {
		return fmt.Errorf("uniform %s not found", variable)
	}

	gl.UniformMatrix4fv(location, 1, false, &mat[0])
	return nil
}

func (s *Shader) UniformTransform(sp gresource.Shader, variable string, transform math.Transform) error {
	return s.UniformMat4(sp, variable, mgl32.Mat4(transform))
}

func (s *Shader) NewShaderLoader() gservice.ShaderLoader {
	return newShaderLoader(s.shaderRepository)
}

func shaderStageToShaderType(stage gresource.ShaderStage) oglresource.ShaderType {
	switch stage {
	case gresource.VertexStage:
		return oglresource.VertexShader
	case gresource.GeometryStage:
		return oglresource.GeometryShader
	case gresource.FragmentStage:
		return oglresource.FragmentShader
	default:
		return 0
	}
}

func getUniformLocation(sp uint32, variable string) int32 {
	if !strings.HasSuffix(variable, "\x00") {
		variable = variable + "\x00"
	}
	return gl.GetUniformLocation(sp, gl.Str(variable))
}
