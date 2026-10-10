package shader

import (
	"context"
	"fmt"
	"game/utils/logger"
	"log/slog"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type Service interface {
	LoadGlslShader(ctx context.Context, code string, stage Stage) (Shader, error)
	BuildProgram(ctx context.Context, builder *ProgramBuilder, name string) (Program, error)
	DeleteProgram(ctx context.Context, program Program) error
	DeleteShader(ctx context.Context, program Shader) error
	ProgramArray() []Program
	ShaderArray() []Shader
}

type service struct {
	shaderRepository  *Repository
	programRepository *ProgramRepository
}

func InitializeShaders(_ context.Context, shaderRepository *Repository, programRepository *ProgramRepository) (Service, error) {
	return &service{
		shaderRepository:  shaderRepository,
		programRepository: programRepository,
	}, nil
}

func (s *service) LoadGlslShader(ctx context.Context, code string, stage Stage) (Shader, error) {
	codeStrPtr, freeCodeStrings := gl.Strs(code)
	defer freeCodeStrings()

	shader := gl.CreateShader(uint32(stage))
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
		return Shader{}, fmt.Errorf("failed to compile shader: %s", string(log))
	}

	shaderObj := Shader{
		ID:    shader,
		Stage: stage,
	}

	s.shaderRepository.Put(shaderObj)

	logger.FromCtx(ctx).LogAttrs(ctx, slog.LevelInfo,
		"loaded glsl shader", logger.Uint32("shader_id", shaderObj.ID))

	return shaderObj, nil
}

func (s *service) BuildProgram(ctx context.Context, builder *ProgramBuilder, name string) (Program, error) {
	vertex, ok := builder.Stage(VertexStage)
	if !ok || vertex.ID <= 0 {
		return Program{}, fmt.Errorf("vertex shader required to create shader program")
	}

	fragment, ok := builder.Stage(FragmentStage)
	if !ok || fragment.ID <= 0 {
		return Program{}, fmt.Errorf("fragment shader required to create shader program")
	}

	geometry, _ := builder.Stage(GeometryStage)

	shaderProgram := gl.CreateProgram()
	gl.AttachShader(shaderProgram, vertex.ID)
	gl.AttachShader(shaderProgram, fragment.ID)
	if geometry.ID > 0 {
		gl.AttachShader(shaderProgram, geometry.ID)
	}
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

		return Program{}, fmt.Errorf("failed to link shader program: %s", string(log))
	}

	program := Program{
		ID:   shaderProgram,
		Name: name,
	}

	if !s.programRepository.Put(program) {
		gl.DeleteProgram(shaderProgram)
		return Program{}, fmt.Errorf("shader program %s already registered", name)
	}

	logger.FromCtx(ctx).LogAttrs(ctx, slog.LevelInfo,
		"loaded shader program",
		logger.Uint32("program_id", program.ID),
		logger.Uint32("vertex_shader", vertex.ID),
		logger.Uint32("fragment_shader", fragment.ID),
		logger.Uint32("geometry_shader", geometry.ID))

	return program, nil
}

func (s *service) DeleteProgram(ctx context.Context, program Program) error {
	var err error

	if gl.IsProgram(program.ID) {
		gl.DeleteProgram(program.ID)
		err = fmt.Errorf("program already deleted or id is not program")
	}

	ok := s.programRepository.Delete(program)
	if !ok {
		logger.FromCtx(ctx).LogAttrs(ctx, slog.LevelWarn, "program deleted from ogl but not found in repository")
	}

	logger.FromCtx(ctx).LogAttrs(ctx, slog.LevelInfo, "shader program deleted", logger.Uint32("program_id", program.ID))
	return err
}

func (s *service) DeleteShader(ctx context.Context, shader Shader) error {
	var err error

	if gl.IsShader(shader.ID) {
		gl.DeleteShader(shader.ID)
		err = fmt.Errorf("shader already deleted or id is not shader")
	}

	ok := s.shaderRepository.Delete(shader)
	if !ok {
		logger.FromCtx(ctx).LogAttrs(ctx, slog.LevelWarn, "shader deleted from ogl but not found in repository")
	}

	return err
}

func (s *service) ProgramArray() []Program {
	return s.programRepository.Array()
}

func (s *service) ShaderArray() []Shader {
	return s.shaderRepository.Array()
}
