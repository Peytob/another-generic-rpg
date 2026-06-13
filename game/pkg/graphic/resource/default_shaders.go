package resource

import (
	"context"
	"game/pkg/utils/logger"
	"log/slog"
)

type Shaders struct {
	World ShaderProgram
}

func (s Shaders) Terminate(ctx context.Context) {
	logger.FromCtx(ctx).Info("deleting shader program", slog.Uint64("shader_program_id", uint64(s.World.Id())))
	s.World.Delete()
}
