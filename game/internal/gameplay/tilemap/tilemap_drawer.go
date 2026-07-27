package tilemap

import (
	"context"
	"engine/graphic/renderer"
	"engine/graphic/resource"
	"engine/math"
	"engine/math/shape"
	"fmt"
	"game/internal/rendering"
)

type Drawer interface {
	Draw(ctx context.Context, tilemap *Tilemap, target renderer.Canvas, opts DrawOpts) error
}

type DrawOpts struct {
	Camera *rendering.Camera
}

type tilemapRenderer struct {
	repository *TileRepository
}

func NewDrawer(tileRepository *TileRepository) Drawer {
	return tilemapRenderer{
		repository: tileRepository,
	}
}

func (t tilemapRenderer) Draw(ctx context.Context, tilemap *Tilemap, target renderer.Canvas, opts DrawOpts) error {
	if tilemap == nil {
		return fmt.Errorf("tilemap is nil")
	}

	if opts.Camera == nil {
		return fmt.Errorf("camera is nil")
	}

	for layerIndex := range tilemap.LayersCount() {
		layer, err := tilemap.Layer(layerIndex)
		if err != nil {
			return fmt.Errorf("unable to get layer %d", layerIndex)
		}

		err = t.drawLayer(ctx, layer, target, opts)
		if err != nil {
			return fmt.Errorf("failed to draw layer %d", layerIndex)
		}
	}

	return nil
}

func (t tilemapRenderer) drawLayer(_ context.Context, layer *Layer, target renderer.Canvas, opts DrawOpts) error {
	transformation := math.NewTransformation()
	sprite := resource.NewSprite(shape.NewRect(tileSize, tileSize), shape.NewZeroRect())

	for x := range layer.Width() {
		for y := range layer.Height() {
			_, err := layer.Tile(x, y)
			if err != nil {
				return fmt.Errorf("unable to get tile (%d, %d)", x, y)
			}

			// todo resolve tile texture here

			transformation.Translate(float32(x)*tileSize, float32(y)*tileSize)

			target.Draw(sprite, &renderer.CanvasOpts{
				Transform: transformation.Transform(),
			})
		}
	}

	return nil
}
