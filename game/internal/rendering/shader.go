package rendering

import (
	"context"
	_ "embed"
	"engine/graphic"
	gresource "engine/graphic/resource"
	"engine/utils/logger"
	"fmt"
)

const (
	TilemapShader = "TilemapShader"
)

//go:embed glsl/src/tilemap.vert
var tilemapVertex string

//go:embed glsl/src/tilemap.frag
var tilemapFragment string

func loadShaders(ctx context.Context, graphic graphic.Graphic) error {
	logger.FromCtx(ctx).Info("loading shaders")

	shaderLoader := graphic.Services().Shader.NewShaderLoader()
	defer shaderLoader.Terminate(ctx)

	tilemapVertexShader, err := shaderLoader.LoadGlslStageShader(tilemapVertex, gresource.VertexStage)
	if err != nil {
		return fmt.Errorf("loading tilemap vertex shader: %w", err)
	}

	tilemapFragmentShader, err := shaderLoader.LoadGlslStageShader(tilemapFragment, gresource.FragmentStage)
	if err != nil {
		return fmt.Errorf("loading tilemap fragment shader: %w", err)
	}

	tilemapShaderDesc := gresource.NewShaderBuilder().
		Set(tilemapVertexShader).
		Set(tilemapFragmentShader)
	_, err = shaderLoader.BuildShaderProgram(tilemapShaderDesc, TilemapShader)
	if err != nil {
		return fmt.Errorf("failed to compile tilemap shader: %w", err)
	}

	logger.FromCtx(ctx).Info("shaders loaded")

	return nil
}
