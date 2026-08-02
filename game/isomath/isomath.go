// Package isomath provides coordinate conversions for 2:1 isometric tile
// rendering in world units.
//
// Coordinate spaces:
//   - Grid (c, r): tile indices (float for entities, int for tile lookup).
//   - Iso  (x, y): continuous isometric space measured in units, where one
//     unit equals one tile width (Unit pixels at render time).
//
// Tile diamonds are TileWu x TileHu units; for a 2:1 dimetric projection
// TileHu = TileWu / 2. GridToScreen returns the position of the diamond's top
// vertex, the conventional anchor for a diamond sprite whose top vertex sits at
// the horizontal center of its bounding box.
//
// The rendering pipeline (camera, projection, sprite quads) works in pixels, so
// use ToPixels/ToUnits to bridge iso-space and pixel-space at the boundaries.
package isomath

import "github.com/go-gl/mathgl/mgl32"

// Unit is the number of pixels per world unit: 1 unit == one tile width == 64 px.
const Unit float32 = 64

// TileWu and TileHu are the tile diamond dimensions in units.
const (
	TileWu float32 = 1
	TileHu float32 = 0.5
)

// GridToScreen maps grid (tile) coordinates to isometric-space coordinates in
// units. The returned point is the top vertex of the tile diamond.
func GridToScreen(c, r float32) mgl32.Vec2 {
	return mgl32.Vec2{
		(c - r) * TileWu / 2,
		(c + r) * TileHu / 2,
	}
}

// ScreenToGrid is the inverse of GridToScreen: it maps an isometric-space point
// (in units) back to the grid coordinates of the tile whose top vertex is there.
func ScreenToGrid(sx, sy float32) mgl32.Vec2 {
	return mgl32.Vec2{
		(sx/(TileWu/2) + sy/(TileHu/2)) / 2,
		(sy/(TileHu/2) - sx/(TileWu/2)) / 2,
	}
}

// ToPixels converts an isometric-space position (units) to pixels.
func ToPixels(v mgl32.Vec2) mgl32.Vec2 {
	return mgl32.Vec2{v.X() * Unit, v.Y() * Unit}
}

// ToUnits converts a pixel position to isometric-space units.
func ToUnits(px mgl32.Vec2) mgl32.Vec2 {
	return mgl32.Vec2{px.X() / Unit, px.Y() / Unit}
}
