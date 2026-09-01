package gamemachine

import (
	"context"
	"engine/hid"
	"fmt"
	"game/internal/gameplay/world"
	"game/internal/gamestate/event"
	"game/internal/gamestate/gameloop"
	"game/internal/gamestate/repositories"
	"game/internal/input"
	"game/internal/rendering"
	"game/internal/sync"
	"time"
)

const PlayingStateIdentifier = StateIdentifier("playing")

// playingLoopTick is the fixed simulation step of the playing state loop.
const playingLoopTick = time.Second / 60

// playingState describes main game state that used while game is running.
// Its Update is a gameloop pipeline: input snapshot and server sync run once
// per frame, simulation advances in fixed ticks, rendering closes the frame.
type playingState struct {
	renderer       rendering.Rendering
	tileRepository *repositories.TileRepository
	hid            hid.Hid
	collector      *input.InputCollector
	loop           *gameloop.Loop
}

func NewPlayingState(r rendering.Rendering, repo repositories.Repositories, h hid.Hid) State {
	s := &playingState{
		renderer:       r,
		tileRepository: repo.TileRepository,
		hid:            h,
		collector:      input.NewInputCollector(),
	}

	s.loop = gameloop.NewLoop(playingLoopTick).
		Before(
			gameloop.InputStage(s.collector),
			gameloop.SyncStage(sync.NoopSync{}),
			gameloop.GlobalInput(),
		).
		EachTick(gameloop.Simulate()).
		After(gameloop.WorldRender(s.renderer))

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

func (s *playingState) Update(ctx context.Context, w *world.World, dt time.Duration) (event.Event, error) {
	return s.loop.Run(ctx, w, dt)
}

// keyboardBindings binds gameplay keys to collector events. Bindings are
// scancode-based, so they follow the current keyboard layout.
func (s *playingState) keyboardBindings() (hid.KeyboardBindings, error) {
	kb := hid.NewKeyboardBindings("playing_keyboard")
	mapper := s.hid.Keyboard()

	bind := func(key hid.Key, action input.Action) error {
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

	for key, action := range map[hid.Key]input.Action{
		hid.KeyW:      input.MoveUp,
		hid.KeyS:      input.MoveDown,
		hid.KeyA:      input.MoveLeft,
		hid.KeyD:      input.MoveRight,
		hid.KeyEscape: input.Exit,
	} {
		if err := bind(key, action); err != nil {
			return hid.KeyboardBindings{}, err
		}
	}

	return kb, nil
}
