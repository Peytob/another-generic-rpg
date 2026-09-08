package loop

import (
	"context"
	"fmt"
	"game/internal/game/state"
	"game/internal/input"
	"game/internal/sync"
	"time"
)

const (
	// maxFrameDt caps the reported frame delta so that long stalls (window
	// drag, debugger pause) do not cause huge catch-up bursts.
	maxFrameDt = 250 * time.Millisecond

	// maxTicksPerFrame caps the number of simulation ticks executed in a
	// single frame; when the cap is hit, accumulated lag is dropped to
	// avoid the spiral of death.
	maxTicksPerFrame = 5
)

// Frame is the shared per-frame context passed to every stage.
type Frame struct {
	// Dt is the real duration of the current frame (already clamped).
	// During EachTick stages it is temporarily set to the fixed tick.
	Dt time.Duration

	// Alpha is the fixed-step interpolation factor in [0, 1): how far the
	// render moment lies between the last two simulation ticks.
	Alpha float64

	// Input is the input snapshot drained by the input stage; it stays
	// empty unless an input stage filled it.
	Input input.InputFrame

	// ServerMessages holds current working server messages
	ServerMessages []sync.ServerMessage

	// ClientMessage holds current working client messages
	ClientMessage []sync.ClientMessage

	State *state.State

	transition state.Event
}

// RequestTransition asks the state machine to perform a transition after
// the current frame completes. The loop never transitions mid-frame.
func (f *Frame) RequestTransition(ev state.Event) {
	f.transition = ev
}

// Transition returns the requested transition or NoEvent when no stage
// requested one.
func (f *Frame) Transition() state.Event {
	return f.transition
}

type Stage func(ctx context.Context, f *Frame) error

// Loop orchestrates the frame pipeline of a single game state.
//
// Loop is not safe for concurrent use; it must be driven from the same
// goroutine that runs the hid Dispatch (the main game loop goroutine).
type Loop struct {
	tick     time.Duration
	before   []Stage
	eachTick []Stage
	after    []Stage
	acc      time.Duration
}

// NewLoop creates a Loop with the given fixed simulation step. A zero or
// negative tick switches the loop into variable-step mode where EachTick
// stages run once per frame with the real frame delta.
func NewLoop(tick time.Duration) *Loop {
	return &Loop{tick: tick}
}

// Before appends stages that run once per frame before simulation ticks.
func (l *Loop) Before(stages ...Stage) *Loop {
	l.before = append(l.before, stages...)
	return l
}

// EachTick appends simulation stages executed with the fixed time step.
func (l *Loop) EachTick(stages ...Stage) *Loop {
	l.eachTick = append(l.eachTick, stages...)
	return l
}

// After appends stages that run once per frame after simulation ticks.
func (l *Loop) After(stages ...Stage) *Loop {
	l.after = append(l.after, stages...)
	return l
}

// Run advances the loop by one frame with the given world and frame delta.
// It returns the transition requested by stages via Frame.RequestTransition
// or NoEvent when no transition was requested.
func (l *Loop) Run(ctx context.Context, currentState *state.State, dt time.Duration) (state.Event, error) {
	if dt > maxFrameDt {
		dt = maxFrameDt
	}

	f := &Frame{Dt: dt, State: currentState}

	if err := l.runStages(ctx, f, l.before); err != nil {
		return state.NoEvent, fmt.Errorf("before stage: %w", err)
	}

	if err := l.advanceTicks(ctx, f); err != nil {
		return state.NoEvent, err
	}

	if err := l.runStages(ctx, f, l.after); err != nil {
		return state.NoEvent, fmt.Errorf("after stage: %w", err)
	}

	return f.transition, nil
}

// advanceTicks executes EachTick stages according to the accumulated time:
// fixed-step loops run them zero..maxTicksPerFrame times per frame, while
// variable-step loops run them exactly once with the real frame delta.
func (l *Loop) advanceTicks(ctx context.Context, f *Frame) error {
	if l.tick <= 0 {
		return l.runTick(ctx, f, f.Dt)
	}

	l.acc += f.Dt

	ticks := 0
	for l.acc >= l.tick && ticks < maxTicksPerFrame {
		if err := l.runTick(ctx, f, l.tick); err != nil {
			return err
		}
		l.acc -= l.tick
		ticks++
	}

	// Still behind after the tick cap: drop the remaining lag instead of
	// trying to catch up forever (spiral of death).
	if ticks == maxTicksPerFrame && l.acc >= l.tick {
		l.acc = 0
	}

	f.Alpha = float64(l.acc) / float64(l.tick)

	return nil
}

// runTick runs all EachTick stages once, exposing dt as the frame delta for
// their duration. The frame is shared, so transition requests made inside
// tick stages are preserved.
func (l *Loop) runTick(ctx context.Context, f *Frame, dt time.Duration) error {
	frameDt := f.Dt
	f.Dt = dt
	defer func() { f.Dt = frameDt }()

	if err := l.runStages(ctx, f, l.eachTick); err != nil {
		return fmt.Errorf("tick stage: %w", err)
	}

	return nil
}

func (l *Loop) runStages(ctx context.Context, f *Frame, stages []Stage) error {
	for _, stage := range stages {
		if err := stage(ctx, f); err != nil {
			return err
		}
	}
	return nil
}
