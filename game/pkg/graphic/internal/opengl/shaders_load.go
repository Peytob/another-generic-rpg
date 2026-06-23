package opengl

import (
	"context"
	_ "embed"
	"fmt"
	"game/pkg/graphic/internal/opengl/resource"
	gresource "game/pkg/graphic/resource"
	"game/pkg/utils/logger"
	"log/slog"

	"github.com/go-gl/gl/v3.3-core/gl"
)

//go:embed glsl/src/world.vert
var worldVertex string

//go:embed glsl/src/world.frag
var worldFragment string

func loadShaders(ctx context.Context) (gresource.Shaders, error) {
	var err error

	logger.FromCtx(ctx).Info("loading shaders")

	/* Shaders loading */

	worldVertexShader, err := buildShader(ctx, worldVertex, resource.VertexShader)
	if err != nil {
		return gresource.Shaders{}, err
	}
	defer deleteShader(ctx, worldVertexShader)

	worldFragmentShader, err := buildShader(ctx, worldFragment, resource.FragmentShader)
	if err != nil {
		return gresource.Shaders{}, err
	}
	defer deleteShader(ctx, worldFragmentShader)

	/* Shader programs loading */

	shaders := gresource.Shaders{}

	shaders.World, err = buildShaderProgram(ctx, resource.ShaderProgramBuilder{
		Vertex:   worldVertexShader,
		Fragment: worldFragmentShader,
	})
	if err != nil {
		return shaders, fmt.Errorf("building shader program failed: %w", err)
	}

	logger.FromCtx(ctx).Info("shaders loaded")

	return shaders, nil
}

func buildShader(ctx context.Context, code string, shaderType resource.ShaderType) (resource.Shader, error) {
	shader, err := resource.BuildShader(ctx, code, shaderType)
	if err != nil {
		return 0, fmt.Errorf("failed to build shader: %w", err)
	}
	logger.FromCtx(ctx).Info("shader build", slog.Uint64("shader_id", uint64(shader.Id())))
	return shader, nil
}

func buildShaderProgram(ctx context.Context, builder resource.ShaderProgramBuilder) (gresource.ShaderProgram, error) {
	shaderProgram, err := resource.BuildShaderProgram(ctx, builder)
	if err != nil {
		return 0, fmt.Errorf("failed to build shader program: %w", err)
	}
	logger.FromCtx(ctx).Info("shader program build", slog.Uint64("shader_program_id", uint64(shaderProgram.Id())))
	return shaderProgram, nil
}

func terminateShaders(ctx context.Context, s gresource.Shaders) {
	logger.FromCtx(ctx).Info("deleting shader program", slog.Uint64("shader_program_id", uint64(s.World.Id())))
	gl.DeleteProgram(s.World.Id())
}

func deleteShader(ctx context.Context, shader resource.Shader) {
	logger.FromCtx(ctx).Info("deleting shader", slog.Uint64("shader_id", uint64(shader.Id())))
	gl.DeleteShader(shader.Id())
}
