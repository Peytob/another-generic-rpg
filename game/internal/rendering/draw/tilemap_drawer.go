package draw

import (
	"context"
	"engine/graphic/renderer"
	"engine/graphic/resource"
	"engine/math"
	"engine/math/shape"
	"fmt"
	"game/internal/gameplay/tilemap"
	"game/internal/rendering"
	"game/isomath"
	stdmath "math"

	"github.com/go-gl/mathgl/mgl32"
)

type TilemapDrawer interface {
	Draw(ctx context.Context, tilemap *tilemap.Tilemap, target renderer.Canvas, opts Opts) error
}

type Opts struct {
	Camera *rendering.Camera
}

type tilemapRenderer struct {
	repository *tilemap.TileRepository
}

func NewTilemapDrawer(tileRepository *tilemap.TileRepository) TilemapDrawer {
	return tilemapRenderer{
		repository: tileRepository,
	}
}

func (t tilemapRenderer) Draw(ctx context.Context, tilemap *tilemap.Tilemap, target renderer.Canvas, opts Opts) error {
	if tilemap == nil {
		return fmt.Errorf("tilemap is nil")
	}

	if opts.Camera == nil {
		return fmt.Errorf("camera is nil")
	}

	for layerIndex := range tilemap.LayersCount() {
		layer, err := tilemap.Layer(layerIndex)
		if err != nil {
			return fmt.Errorf("unable to get layer %d: %w", layerIndex, err)
		}

		err = t.drawLayer(ctx, layer, target, opts)
		if err != nil {
			return fmt.Errorf("failed to draw layer %d: %w", layerIndex, err)
		}
	}

	return nil
}

func (t tilemapRenderer) drawLayer(_ context.Context, layer *tilemap.Layer, target renderer.Canvas, opts Opts) error {
	transformation := math.NewTransformation()
	sprite := resource.NewSprite(shape.NewRect(tileWPx, tileHPx), shape.NewZeroRect())

	startX, startY, endX, endY := visibleTileRange(layer, opts.Camera)

	// todo iso depth sorting: once tall sprites/objects or overlapping layers
	//      are drawn, render back-to-front ordered by screen Y
	//      (isomath.GridToScreen(x, y).Y(), i.e. (x+y) ascending) so nearer
	//      tiles correctly occlude farther ones.
	for x := startX; x < endX; x++ {
		for y := startY; y < endY; y++ {
			_, err := layer.Tile(x, y)
			if err != nil {
				return fmt.Errorf("unable to get tile (%d, %d): %w", x, y, err)
			}

			// todo resolve tile texture here

			pos := isomath.ToPixels(isomath.GridToScreen(float32(x), float32(y)))
			transformation.Translate(pos.X()-tileWPx/2, pos.Y())

			target.Draw(sprite, &renderer.CanvasOpts{
				Transform: transformation.Transform(),
			})
		}
	}

	return nil
}

func visibleTileRange(layer *tilemap.Layer, camera *rendering.Camera) (startX, startY, endX, endY int) {
	camPos := camera.GetPosition()
	camArea := camera.GetArea()

	left := camPos.X() - tileWPx/2
	right := camPos.X() + camArea.X() + tileWPx/2
	top := camPos.Y() - tileHPx
	bottom := camPos.Y() + camArea.Y() + tileHPx

	corners := [4]mgl32.Vec2{
		{left, top},
		{right, top},
		{left, bottom},
		{right, bottom},
	}

	minC, maxC := float32(stdmath.MaxInt32), float32(-stdmath.MaxInt32)
	minR, maxR := float32(stdmath.MaxInt32), float32(-stdmath.MaxInt32)
	for _, corner := range corners {
		u := isomath.ToUnits(corner)
		g := isomath.ScreenToGrid(u.X(), u.Y())
		minC = min(minC, g.X())
		maxC = max(maxC, g.X())
		minR = min(minR, g.Y())
		maxR = max(maxR, g.Y())
	}

	const margin = 1
	startX = clampInt(floorToInt(minC)-margin, 0, layer.Width())
	endX = clampInt(ceilToInt(maxC)+margin, 0, layer.Width())
	startY = clampInt(floorToInt(minR)-margin, 0, layer.Height())
	endY = clampInt(ceilToInt(maxR)+margin, 0, layer.Height())
	return startX, startY, endX, endY
}

func clampInt(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func floorToInt(x float32) int {
	i := int(x)
	if x < 0 && float32(i) != x {
		i--
	}
	return i
}

func ceilToInt(x float32) int {
	i := int(x)
	if x > 0 && float32(i) != x {
		i++
	}
	return i
}
