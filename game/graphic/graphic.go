package graphic

import (
	"context"
	"errors"
	"fmt"
	"game/graphic/shader"
	"game/graphic/telemetry"
	"game/graphic/textures"
	"game/graphic/uniform"
	"game/pkg/files"
	"game/pkg/logger"
	"log/slog"
)

type repositories struct {
	shader       *shader.Repository
	program      *shader.ProgramRepository
	uniformBlock *uniform.BlockRepository
	texture      *textures.Repository
}

type Graphic struct {
	repositories repositories

	telemetry telemetry.Telemetry
	shader    shader.Service
	uniform   uniform.Service
	textures  textures.Service
}

type Opts struct {
}

func Initialize(ctx context.Context, _ Opts) (*Graphic, error) {
	ctx = logger.WithLogValue(ctx, slog.String("svc", "graphic_initialization"))
	var err error
	g := &Graphic{}

	g.repositories = repositories{}
	g.repositories.shader = shader.NewShaderRepository()
	g.repositories.program = shader.NewProgramRepository()
	g.repositories.uniformBlock = uniform.NewUniformBlockRepository()
	g.repositories.texture = textures.NewTextureRepository()

	g.telemetry, err = telemetry.InitializeTelemetry(ctx)
	if err != nil {
		return nil, fmt.Errorf("telemetry initialize: %w", err)
	}

	g.shader, err = shader.InitializeShaders(ctx, g.repositories.shader, g.repositories.program)
	if err != nil {
		return nil, fmt.Errorf("shader initialization: %w", err)
	}

	g.uniform, err = uniform.InitializeUniforms(ctx, g.repositories.uniformBlock)
	if err != nil {
		return nil, fmt.Errorf("uniform initialization: %w", err)
	}

	g.textures, err = textures.InitializeTextures(ctx, files.NewStubLoader(), g.repositories.texture)
	if err != nil {
		return nil, fmt.Errorf("textures initialization: %w", err)
	}

	return g, nil
}

func (g *Graphic) Telemetry() telemetry.Telemetry {
	return g.telemetry
}

func (g *Graphic) Terminate(ctx context.Context) error {
	ctx = logger.WithLogValue(ctx, slog.String("svc", "graphic_termination"))
	var joinedErr error

	logger.FromCtx(ctx).LogAttrs(ctx, slog.LevelInfo, "terminating graphic")

	logger.FromCtx(ctx).LogAttrs(ctx, slog.LevelInfo, "terminating shader programs")
	for _, program := range g.shader.ProgramArray() {
		err := g.shader.DeleteProgram(ctx, program)
		if err != nil {
			logger.FromCtx(ctx).LogAttrs(ctx, slog.LevelError, "shader program delete", logger.Error(err))
			joinedErr = errors.Join(joinedErr, err)
		}
	}

	logger.FromCtx(ctx).LogAttrs(ctx, slog.LevelInfo, "terminating shaders")
	for _, shaderObj := range g.shader.ShaderArray() {
		err := g.shader.DeleteShader(ctx, shaderObj)
		if err != nil {
			logger.FromCtx(ctx).LogAttrs(ctx, slog.LevelError, "shader delete", logger.Error(err))
			joinedErr = errors.Join(joinedErr, err)
		}
	}

	logger.FromCtx(ctx).LogAttrs(ctx, slog.LevelInfo, "terminating uniform blocks")
	for _, block := range g.uniform.BlockArray() {
		err := g.uniform.DeleteBlock(ctx, block)
		if err != nil {
			logger.FromCtx(ctx).LogAttrs(ctx, slog.LevelError, "uniform block delete", logger.Error(err))
			joinedErr = errors.Join(joinedErr, err)
		}
	}

	logger.FromCtx(ctx).LogAttrs(ctx, slog.LevelInfo, "terminating textures")
	for _, texture := range g.textures.TextureArray() {
		err := g.textures.DeleteTexture(ctx, texture)
		if err != nil {
			logger.FromCtx(ctx).LogAttrs(ctx, slog.LevelError, "texture delete", logger.Error(err))
			joinedErr = errors.Join(joinedErr, err)
		}
	}

	return joinedErr
}

func (g *Graphic) Uniform() uniform.Service {
	return g.uniform
}

func (g *Graphic) Textures() textures.Service {
	return g.textures
}
