package gamemachine

import (
	"context"
	"engine/hid"
	"fmt"
	"game/internal/gameplay/world"
	"game/internal/gamestate/gameloop"
	"game/internal/gamestate/repositories"
	"game/internal/rendering"
	"time"
)

const PlayingStateIdentifier = StateIdentifier("playing")

// playingLoopTick is the fixed simulation step of the playing state loop.
const playingLoopTick = time.Second / 60

// gameAction enumerates input actions recognized by the playing state.
type gameAction int

const (
	actionMoveUp gameAction = iota
	actionMoveDown
	actionMoveLeft
	actionMoveRight
	actionExit
)

// playingState describes main game state that used while game is running.
// Its Update is a gameloop pipeline: input snapshot and server sync run once
// per frame, simulation advances in fixed ticks, rendering closes the frame.
type playingState struct {
	renderer       rendering.Rendering
	tileRepository *repositories.TileRepository
	hid            hid.Hid
	collector      *gameloop.InputCollector[gameAction]
	loop           *gameloop.Loop[Event, gameAction]
}

func NewPlayingState(r rendering.Rendering, repo repositories.Repositories, h hid.Hid) State {
	s := &playingState{
		renderer:       r,
		tileRepository: repo.TileRepository,
		hid:            h,
		collector:      gameloop.NewInputCollector[gameAction](),
	}

	s.loop = gameloop.NewLoop[Event, gameAction](playingLoopTick).
		Before(
			gameloop.InputStage[Event, gameAction](s.collector),
			gameloop.SyncStage[Event, gameAction](gameloop.NoopSync{}),
			s.handleGlobalInput,
		).
		EachTick(s.simulate).
		After(s.render)

	return s
}

func (s *playingState) Identifier() StateIdentifier {
	return PlayingStateIdentifier
}

func (s *playingState) OnEnter(_ context.Context, _ *world.World) error {
	kb, err := s.keyboardBindings()
	if err != nil {
		return fmt.Errorf("failed to build keyboard bindings: %w", err)
	}

	s.hid.Keyboard().SetCurrentBindings(kb)
	s.hid.Mouse().SetCurrentBindings(hid.NewMouseBindings("playing_mouse"))

	return nil
}

func (s *playingState) OnExit(_ context.Context, _ *world.World) error {
	// Reset bindings so callbacks of a deactivated state never fire.
	s.hid.Keyboard().SetCurrentBindings(hid.NewKeyboardBindings("unbound_keyboard"))
	s.hid.Mouse().SetCurrentBindings(hid.NewMouseBindings("unbound_mouse"))

	return nil
}

func (s *playingState) Update(ctx context.Context, w *world.World, dt time.Duration) (Event, error) {
	return s.loop.Run(ctx, w, dt)
}

// handleGlobalInput maps frame input to state machine transitions.
func (s *playingState) handleGlobalInput(_ context.Context, f *gameloop.Frame[Event, gameAction]) error {
	if f.Input.Pressed(actionExit) {
		f.RequestTransition(StoppedEvent)
	}
	return nil
}

// simulate advances the gameplay simulation by one fixed tick.
func (s *playingState) simulate(_ context.Context, _ *gameloop.Frame[Event, gameAction]) error {
	// todo advance world simulation by f.Dt applying f.ServerMessages
	return nil
}

// render draws the current frame.
func (s *playingState) render(_ context.Context, _ *gameloop.Frame[Event, gameAction]) error {
	// todo draw scene via s.renderer.Drawers() using f.Alpha interpolation
	return nil
}

// keyboardBindings binds gameplay keys to collector events. Bindings are
// scancode-based, so they follow the current keyboard layout.
func (s *playingState) keyboardBindings() (hid.KeyboardBindings, error) {
	kb := hid.NewKeyboardBindings("playing_keyboard")
	mapper := s.hid.Keyboard()

	bind := func(key hid.Key, action gameAction) error {
		scancode := mapper.GetScancode(key)

		press := func(_ hid.Key, _ int, _ hid.Action, _ hid.Modifier) {
			s.collector.Press(action)
		}
		release := func(_ hid.Key, _ int, _ hid.Action, _ hid.Modifier) {
			s.collector.Release(action)
		}

		if err := kb.AddBinding(hid.KeyboardBinding{Scancode: scancode, Action: hid.Pressed, Mods: hid.ModNone}, press); err != nil {
			return fmt.Errorf("bind key %d press: %w", key, err)
		}
		if err := kb.AddBinding(hid.KeyboardBinding{Scancode: scancode, Action: hid.Released, Mods: hid.ModNone}, release); err != nil {
			return fmt.Errorf("bind key %d release: %w", key, err)
		}

		return nil
	}

	for key, action := range map[hid.Key]gameAction{
		hid.KeyW:      actionMoveUp,
		hid.KeyS:      actionMoveDown,
		hid.KeyA:      actionMoveLeft,
		hid.KeyD:      actionMoveRight,
		hid.KeyEscape: actionExit,
	} {
		if err := bind(key, action); err != nil {
			return hid.KeyboardBindings{}, err
		}
	}

	return kb, nil
}
