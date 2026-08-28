package ecs

import (
	"context"
	"errors"
	"testing"
)

func TestSubscribeAndEmit(t *testing.T) {
	t.Parallel()

	t.Run("delivers event to single subscriber", func(t *testing.T) {
		t.Parallel()

		w := NewWorld()
		var received *DamageEvent

		Subscribe[*DamageEvent](w, func(_ context.Context, e *DamageEvent) error {
			received = e
			return nil
		})

		if err := w.EmitEvent(context.Background(), &DamageEvent{Amount: 42}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if received == nil {
			t.Fatal("event was not delivered")
		}
		assertEqual(t, received.Amount, 42)
	})

	t.Run("delivers event to multiple subscribers in order", func(t *testing.T) {
		t.Parallel()

		w := NewWorld()
		var order []int

		Subscribe[*DamageEvent](w, func(context.Context, *DamageEvent) error {
			order = append(order, 1)
			return nil
		})
		Subscribe[*DamageEvent](w, func(context.Context, *DamageEvent) error {
			order = append(order, 2)
			return nil
		})
		Subscribe[*DamageEvent](w, func(context.Context, *DamageEvent) error {
			order = append(order, 3)
			return nil
		})

		if err := w.EmitEvent(context.Background(), &DamageEvent{Amount: 1}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		assertSliceEqual(t, order, []int{1, 2, 3})
	})

	t.Run("different event types are isolated", func(t *testing.T) {
		t.Parallel()

		w := NewWorld()
		damageCount := 0
		healCount := 0

		Subscribe[*DamageEvent](w, func(context.Context, *DamageEvent) error {
			damageCount++
			return nil
		})
		Subscribe[*HealEvent](w, func(context.Context, *HealEvent) error {
			healCount++
			return nil
		})

		_ = w.EmitEvent(context.Background(), &DamageEvent{Amount: 10})
		_ = w.EmitEvent(context.Background(), &HealEvent{Amount: 5})

		assertEqual(t, damageCount, 1)
		assertEqual(t, healCount, 1)
	})

	t.Run("emit with no subscribers returns nil", func(t *testing.T) {
		t.Parallel()

		w := NewWorld()

		if err := w.EmitEvent(context.Background(), &DamageEvent{Amount: 1}); err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
	})
}

func TestUnsubscribe(t *testing.T) {
	t.Parallel()

	t.Run("stops delivery after unsubscribe", func(t *testing.T) {
		t.Parallel()

		w := NewWorld()
		callCount := 0

		sub := Subscribe[*DamageEvent](w, func(context.Context, *DamageEvent) error {
			callCount++
			return nil
		})

		_ = w.EmitEvent(context.Background(), &DamageEvent{Amount: 1})
		assertEqual(t, callCount, 1)

		sub.Unsubscribe()

		_ = w.EmitEvent(context.Background(), &DamageEvent{Amount: 2})
		assertEqual(t, callCount, 1)
	})

	t.Run("only unsubscribed handler stops, others continue", func(t *testing.T) {
		t.Parallel()

		w := NewWorld()
		calls := []string{}

		sub1 := Subscribe[*DamageEvent](w, func(context.Context, *DamageEvent) error {
			calls = append(calls, "a")
			return nil
		})
		Subscribe[*DamageEvent](w, func(context.Context, *DamageEvent) error {
			calls = append(calls, "b")
			return nil
		})

		sub1.Unsubscribe()

		_ = w.EmitEvent(context.Background(), &DamageEvent{Amount: 1})

		assertSliceEqual(t, calls, []string{"b"})
	})

	t.Run("unsubscribe all then emit cleans up and returns nil", func(t *testing.T) {
		t.Parallel()

		w := NewWorld()
		sub := Subscribe[*DamageEvent](w, func(context.Context, *DamageEvent) error {
			return nil
		})

		sub.Unsubscribe()

		if err := w.EmitEvent(context.Background(), &DamageEvent{Amount: 1}); err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})
}

func TestEmitNilEvent(t *testing.T) {
	t.Parallel()

	t.Run("emit nil event is a no-op", func(t *testing.T) {
		t.Parallel()

		w := NewWorld()
		called := false

		Subscribe[*DamageEvent](w, func(context.Context, *DamageEvent) error {
			called = true
			return nil
		})

		if err := w.EmitEvent(context.Background(), nil); err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
		if called {
			t.Error("handler should not be called for nil event")
		}
	})
}

func TestEmitEventError(t *testing.T) {
	t.Parallel()

	t.Run("stops on first error and returns it", func(t *testing.T) {
		t.Parallel()

		w := NewWorld()
		errSentinel := errors.New("handler failed")
		called := []string{}

		Subscribe[*DamageEvent](w, func(context.Context, *DamageEvent) error {
			called = append(called, "first")
			return nil
		})
		Subscribe[*DamageEvent](w, func(context.Context, *DamageEvent) error {
			called = append(called, "second")
			return errSentinel
		})
		Subscribe[*DamageEvent](w, func(context.Context, *DamageEvent) error {
			called = append(called, "third")
			return nil
		})

		err := w.EmitEvent(context.Background(), &DamageEvent{Amount: 1})

		if !errors.Is(err, errSentinel) {
			t.Errorf("expected errSentinel, got %v", err)
		}
		assertSliceEqual(t, called, []string{"first", "second"})
	})
}

func TestEmitEventReentrancy(t *testing.T) {
	t.Parallel()

	t.Run("unsubscribe during emit does not panic", func(t *testing.T) {
		t.Parallel()

		w := NewWorld()

		var sub2 Subscription
		sub2 = Subscribe[*DamageEvent](w, func(context.Context, *DamageEvent) error {
			sub2.Unsubscribe()
			return nil
		})

		Subscribe[*DamageEvent](w, func(context.Context, *DamageEvent) error {
			return nil
		})

		if err := w.EmitEvent(context.Background(), &DamageEvent{Amount: 1}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("emit during emit is safe", func(t *testing.T) {
		t.Parallel()

		w := NewWorld()
		emitCount := 0

		Subscribe[*DamageEvent](w, func(ctx context.Context, _ *DamageEvent) error {
			emitCount++
			if emitCount == 1 {
				_ = w.EmitEvent(ctx, &DamageEvent{Amount: 2})
			}
			return nil
		})

		if err := w.EmitEvent(context.Background(), &DamageEvent{Amount: 1}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		assertEqual(t, emitCount, 2)
	})
}
