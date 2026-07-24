package tilemap

import "errors"

// ErrOutOfRange requested index lies outside the valid bounds.
var ErrOutOfRange = errors.New("index out of range")
