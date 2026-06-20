package graphic

import (
	"context"
	_ "embed"
	"fmt"
	gresource "game/pkg/graphic/resource"
	"game/pkg/utils/logger"
	"log/slog"
)

//go:embed glsl/src/world.vert
var worldVertex string

//go:embed glsl/src/world.frag
var worldFragment string

func LoadShaders(ctx context.Context) (*gresource.Shaders, error) {
	shaders := &gresource.Shaders{}
	var err error

	logger.FromCtx(ctx).Info("loading shaders")

	/* Shaders loading */

	worldVertexShader, err := buildShader(ctx, worldVertex, gresource.VertexShader)
	if err != nil {
		return nil, err
	}
	defer deleteShader(ctx, worldVertexShader)

	worldFragmentShader, err := buildShader(ctx, worldFragment, gresource.FragmentShader)
	if err != nil {
		return nil, err
	}
	defer deleteShader(ctx, worldFragmentShader)

	/* Shader programs loading */

	shaders.World, err = buildShaderProgram(ctx, gresource.ShaderProgramBuilder{
		Vertex:   worldVertexShader,
		Fragment: worldFragmentShader,
	})
	if err != nil {
		return nil, fmt.Errorf("building shader program failed: %w", err)
	}

	logger.FromCtx(ctx).Info("shaders loaded")

	return shaders, nil
}

func buildShader(ctx context.Context, code string, shaderType gresource.ShaderType) (gresource.Shader, error) {
	shader, err := gresource.BuildShader(ctx, code, shaderType)
	if err != nil {
		return 0, fmt.Errorf("failed to build shader: %w", err)
	}
	logger.FromCtx(ctx).Info("shader build", slog.Uint64("shader_id", uint64(shader.Id())))
	return shader, nil
}

func buildShaderProgram(ctx context.Context, builder gresource.ShaderProgramBuilder) (gresource.ShaderProgram, error) {
	shaderProgram, err := gresource.BuildShaderProgram(ctx, builder)
	if err != nil {
		return 0, fmt.Errorf("failed to build shader program: %w", err)
	}
	logger.FromCtx(ctx).Info("shader program build", slog.Uint64("shader_program_id", uint64(shaderProgram.Id())))
	return shaderProgram, nil
}

func deleteShader(ctx context.Context, shader gresource.Shader) {
	logger.FromCtx(ctx).Info("deleting shader", slog.Uint64("shader_id", uint64(shader.Id())))
	shader.Delete()
}
