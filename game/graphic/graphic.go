package graphic

import (
	"context"
	"game/graphic/shader"
	"game/graphic/telemetry"
	"game/graphic/uniform"
	"game/utils/logger"
	"log/slog"
)

type repositories struct {
	shader       *shader.Repository
	program      *shader.ProgramRepository
	uniformBlock *uniform.BlockRepository
}

type Graphic struct {
	repositories repositories

	telemetry telemetry.Telemetry
	shader    shader.Service
	uniform   uniform.Service
}

type Opts struct {
}

func Init(_ Opts) (*Graphic, error) {
	g := &Graphic{}

	g.repositories = repositories{}
	g.repositories.shader = shader.NewShaderRepository()
	g.repositories.program = shader.NewProgramRepository()
	g.repositories.uniformBlock = uniform.NewUniformBlockRepository()

	g.telemetry = telemetry.NewTelemetry()
	g.shader = shader.NewShaderService(g.repositories.shader, g.repositories.program)
	g.uniform = uniform.NewUniformService(g.repositories.uniformBlock)

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
		}
	}

	logger.FromCtx(ctx).LogAttrs(ctx, slog.LevelInfo, "terminating shaders")
	for _, shaderObj := range g.shader.ShaderArray() {
		err := g.shader.DeleteShader(ctx, shaderObj)
		if err != nil {
			logger.FromCtx(ctx).LogAttrs(ctx, slog.LevelError, "shader delete", logger.Error(err))
		}
	}

	return joinedErr
}
