package loop

import "context"

type Step func(ctx context.Context) error

type Loop struct {
	steps []Step
}

type Builder struct {
	steps []Step
}
