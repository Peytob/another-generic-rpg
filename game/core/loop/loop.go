package loop

import "context"

type Step func(ctx context.Context) error

type Loop struct{}

type Builder struct{}
