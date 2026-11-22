package files

import "context"

type Loader interface {
	LoadImage(ctx context.Context, name string) (Image, error)
}

type stubLoader struct{}

func NewStubLoader() Loader {
	return &stubLoader{}
}

func (stubLoader) LoadImage(_ context.Context, _ string) (Image, error) {
	return Image{
		Width:  1,
		Height: 1,
		Pixels: []byte{0xFF, 0xFF, 0xFF, 0xFF},
	}, nil
}
