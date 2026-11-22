package machine

import (
	"game/input"
	"time"
)

type Frame struct {
	// Dt is the real duration of the current frame (already clamped).
	// During EachTick stages it is temporarily set to the fixed tick.
	Dt time.Duration

	// Input is the input snapshot drained by the input module
	Input input.Frame
}
