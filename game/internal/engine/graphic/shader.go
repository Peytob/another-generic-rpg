package graphic

import (
	"context"
	_ "embed"
	"fmt"
	"game/pkg/graphic"
	gresource "game/pkg/graphic/resource"
	"game/pkg/utils/logger"
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
		return fmt.Errorf("loading tilemap vertex shader: %w", err)
	}

	tilemapShaderDesc := gresource.NewShaderBuilder().
		Set(tilemapVertexShader).
		Set(tilemapFragmentShader)
	tilemapShader, err := shaderLoader.BuildShaderProgram(tilemapShaderDesc, TilemapShader)
	if err != nil {
		return fmt.Errorf("failed to compile tilemap shader: %w", err)
	}
	graphic.Repositories().Shader.Put(tilemapShader)

	logger.FromCtx(ctx).Info("shaders loaded")

	return nil
}
