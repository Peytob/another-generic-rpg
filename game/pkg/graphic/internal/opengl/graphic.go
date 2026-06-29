package opengl

import (
	"context"
	"game/pkg/graphic"
	oglrenderer "game/pkg/graphic/internal/opengl/renderer"
	oglservice "game/pkg/graphic/internal/opengl/service"
	grenderer "game/pkg/graphic/renderer"
	grepository "game/pkg/graphic/repository"
	gservice "game/pkg/graphic/service"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type Graphic struct {
	renderer     *oglrenderer.Renderer
	services     gservice.Services
	repositories grepository.Repositories
}

func NewGraphic(ctx context.Context) (*Graphic, error) {
	repositories := grepository.Repositories{
		Uniform: grepository.NewUniformBlockRepository(),
		Shader:  grepository.NewShaderRepository(),
	}

	services := gservice.Services{}
	services.Shader = oglservice.NewShader(repositories.Shader)
	services.Uniform = oglservice.NewUniformBlock(repositories.Uniform)

	renderer := oglrenderer.NewRenderer(ctx, services.Shader)

	return &Graphic{
		renderer:     renderer,
		services:     services,
		repositories: repositories,
	}, nil
}

func (g Graphic) Renderer() grenderer.Renderer {
	return g.renderer
}

func (g Graphic) NewCanvas() grenderer.Canvas {
	return oglrenderer.NewCanvas()
}

func (g Graphic) Services() gservice.Services {
	return g.services
}

func (g Graphic) Repositories() grepository.Repositories {
	return g.repositories
}

func (g Graphic) Telemetry() graphic.Telemetry {
	return graphic.Telemetry{
		Name:     gl.GoStr(gl.GetString(gl.VENDOR)),
		Version:  gl.GoStr(gl.GetString(gl.VERSION)),
		Renderer: gl.GoStr(gl.GetString(gl.RENDERER)),
	}
}

func (g Graphic) Terminate(ctx context.Context) {
	//terminateShaders(ctx, g.shaders)
	//terminateUniformBlocks(ctx, g.uniformBlocks)
	g.renderer.Terminate(ctx)
}
