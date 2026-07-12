package opengl

import (
	"context"
	"engine/graphic"
	oglrenderer "engine/graphic/internal/opengl/renderer"
	oglservice "engine/graphic/internal/opengl/service"
	grenderer "engine/graphic/renderer"
	grepository "engine/graphic/repository"
	gservice "engine/graphic/service"
	"engine/utils/logger"

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

func (g *Graphic) Renderer() grenderer.Renderer {
	return g.renderer
}

func (g *Graphic) NewCanvas() grenderer.Canvas {
	return oglrenderer.NewCanvas()
}

func (g *Graphic) Services() gservice.Services {
	return g.services
}

func (g *Graphic) Repositories() grepository.Repositories {
	return g.repositories
}

func (g *Graphic) Telemetry() graphic.Telemetry {
	return graphic.Telemetry{
		Name:     gl.GoStr(gl.GetString(gl.VENDOR)),
		Version:  gl.GoStr(gl.GetString(gl.VERSION)),
		Renderer: gl.GoStr(gl.GetString(gl.RENDERER)),
	}
}

func (g *Graphic) Terminate(ctx context.Context) {
	g.terminateShaders(ctx)
	g.terminateUniformBlocks(ctx)
	g.renderer.Terminate(ctx)
}

func (g *Graphic) terminateShaders(ctx context.Context) {
	for _, shader := range g.repositories.Shader.All() {
		logger.FromCtx(ctx).Info("deleting shader program", "id", shader.ID, "name", shader.Name)
		gl.DeleteProgram(shader.ID)
	}
}

func (g *Graphic) terminateUniformBlocks(ctx context.Context) {
	for _, ub := range g.repositories.Uniform.All() {
		logger.FromCtx(ctx).Info("deleting uniform block", "id", ub.ID, "name", ub.Name)
		gl.DeleteBuffers(1, &ub.ID)
	}
}
