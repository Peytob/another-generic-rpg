package tilemap

import "errors"

// ErrOutOfRange requested index lies outside the valid bounds.
var ErrOutOfRange = errors.New("index out of range")

// ErrInvalidDimensions layers count or layer dimensions are invalid.
var ErrInvalidDimensions = errors.New("invalid dimensions")
