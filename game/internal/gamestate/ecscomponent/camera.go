package ecscomponent

import (
	rtypes "game/internal/rendering/types"
)

//go:generate go run typeid
type CameraComponent struct {
	Camera *rtypes.Camera
}

// CameraComponentTypeID is the generated type ID (FNV-1a 64-bit) of
// game/internal/gamestate/ecscomponent.CameraComponent. DO NOT EDIT.
const CameraComponentTypeID int64 = 2833317103761503830
