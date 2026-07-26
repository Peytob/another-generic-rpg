package tilemap

import (
	"context"
	"engine/graphic/renderer"
	"engine/graphic/resource"
	"engine/math"
	"engine/math/shape"
	"fmt"
)

type Drawer interface {
	Draw(ctx context.Context, tilemap *Tilemap, target renderer.Canvas) error
}

type tilemapRenderer struct {
	repository *TileRepository
}

func NewDrawer(tileRepository *TileRepository) Drawer {
	return tilemapRenderer{
		repository: tileRepository,
	}
}

func (t tilemapRenderer) Draw(ctx context.Context, tilemap *Tilemap, target renderer.Canvas) error {
	for layerIndex := range tilemap.LayersCount() {
		layer, err := tilemap.Layer(layerIndex)
		if err != nil {
			return fmt.Errorf("unable to get layer %d", layerIndex)
		}

		err = t.drawLayer(ctx, layer, target)
		if err != nil {
			return fmt.Errorf("failed to draw layer %d", layerIndex)
		}
	}

	return nil
}

func (t tilemapRenderer) drawLayer(ctx context.Context, layer *Layer, target renderer.Canvas) error {
	for x := range layer.Width() {
		for y := range layer.Height() {
			tile, err := layer.Tile(x, y)
			if err != nil {
				return fmt.Errorf("unable to get tile (%d, %d)", x, y)
			}

			//if tile == nil {
			//	continue
			//}

			t.drawTile(ctx, tile, x, y, target)
		}
	}

	return nil
}

func (t tilemapRenderer) drawTile(_ context.Context, _ *Tile, x, y int, target renderer.Canvas) {
	// todo resolve tile texture here

	const tileSize = 32

	transformation := math.NewTransformation()
	transformation.Translate(float32(x)*tileSize, float32(y)*tileSize)

	sprite := resource.NewSprite(shape.NewRect(tileSize, tileSize), shape.NewZeroRect())

	// todo move draw tile to reuse transformations
	target.Draw(sprite, &renderer.CanvasOpts{
		Transform: transformation.Transform(),
	})
}
