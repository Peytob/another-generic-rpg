package drawable

import (
	"context"
	"game/graphic/textures"
)

type VertexDrawer struct {
}

type DrawOpts struct {
	Texture textures.Texture
}

func (d VertexDrawer) Draw(_ context.Context, _ Drawable, _ DrawOpts) {

}
