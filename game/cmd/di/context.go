package di

import "context"

type Context struct {
	context.Context
}

func WrapContext(ctx context.Context) Context {
	return Context{
		Context: ctx,
	}
}
