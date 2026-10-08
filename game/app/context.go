package app

import (
	"context"
	"game/utils/logger"
	"log/slog"
	"time"
)

type Context struct {
	ctx context.Context
}

func WrapCtx(ctx context.Context) Context {
	casted, ok := ctx.(Context)
	if !ok {
		return Context{ctx: ctx}
	}

	return casted
}

func (c Context) Deadline() (deadline time.Time, ok bool) {
	return c.ctx.Deadline()
}

func (c Context) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c Context) Err() error {
	return c.ctx.Err()
}

func (c Context) Value(key any) any {
	return c.ctx.Value(key)
}

func (c Context) WithLogValue(value slog.Attr) Context {
	c.ctx = logger.WithLogValue(c.ctx, value)
	return c
}
