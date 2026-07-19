package ecs

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestAddSystem(t *testing.T) {
	t.Parallel()

	t.Run("Systems returns all registered systems in order", func(t *testing.T) {
		t.Parallel()

		w := NewWorld()
		s1 := func(context.Context, World, time.Duration) error { return nil }
		s2 := func(context.Context, World, time.Duration) error { return nil }
		s3 := func(context.Context, World, time.Duration) error { return nil }

		w.AddSystem(s1)
		w.AddSystem(s2)
		w.AddSystem(s3)

		systems := w.Systems()
		assertEqual(t, len(systems), 3)
	})

	t.Run("Systems returns a copy", func(t *testing.T) {
		t.Parallel()

		w := NewWorld()
		w.AddSystem(func(context.Context, World, time.Duration) error { return nil })

		systems := w.Systems()
		systems[0] = nil

		again := w.Systems()
		if again[0] == nil {
			t.Error("Systems() should return a copy")
		}
	})

	t.Run("SystemsIter iterates in registration order", func(t *testing.T) {
		t.Parallel()

		w := NewWorld()
		w.AddSystem(func(context.Context, World, time.Duration) error { return nil })
		w.AddSystem(func(context.Context, World, time.Duration) error { return nil })

		count := 0
		for range w.SystemsIter() {
			count++
		}
		assertEqual(t, count, 2)
	})
}

func TestUpdate(t *testing.T) {
	t.Parallel()

	t.Run("executes systems in registration order", func(t *testing.T) {
		t.Parallel()

		w := NewWorld()
		var order []int

		w.AddSystem(func(context.Context, World, time.Duration) error {
			order = append(order, 1)
			return nil
		})
		w.AddSystem(func(context.Context, World, time.Duration) error {
			order = append(order, 2)
			return nil
		})
		w.AddSystem(func(context.Context, World, time.Duration) error {
			order = append(order, 3)
			return nil
		})

		if err := w.Update(context.Background(), time.Millisecond); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		assertSliceEqual(t, order, []int{1, 2, 3})
	})

	t.Run("stops on first error", func(t *testing.T) {
		t.Parallel()

		w := NewWorld()
		errSentinel := errors.New("system error")
		executed := []int{}

		w.AddSystem(func(context.Context, World, time.Duration) error {
			executed = append(executed, 1)
			return nil
		})
		w.AddSystem(func(context.Context, World, time.Duration) error {
			executed = append(executed, 2)
			return errSentinel
		})
		w.AddSystem(func(context.Context, World, time.Duration) error {
			executed = append(executed, 3)
			return nil
		})

		err := w.Update(context.Background(), 0)
		if !errors.Is(err, errSentinel) {
			t.Errorf("expected errSentinel, got %v", err)
		}
		assertSliceEqual(t, executed, []int{1, 2})
	})

	t.Run("passes correct dt to systems", func(t *testing.T) {
		t.Parallel()

		w := NewWorld()
		dtWanted := 33 * time.Millisecond

		w.AddSystem(func(_ context.Context, _ World, dt time.Duration) error {
			if dt != dtWanted {
				t.Errorf("dt = %v, want %v", dt, dtWanted)
			}
			return nil
		})

		_ = w.Update(context.Background(), dtWanted)
	})

	t.Run("system can access world components", func(t *testing.T) {
		t.Parallel()

		w := NewWorld()
		e := w.NewEntity()
		w.RegisterComponent(e, Health{HP: 10})

		w.AddSystem(func(_ context.Context, world World, _ time.Duration) error {
			entities := world.Query(hpType)
			for _, ent := range entities {
				hp, _ := GetComponent[Health](world, ent)
				hp.HP -= 5
				world.RegisterComponent(ent, hp)
			}
			return nil
		})

		_ = w.Update(context.Background(), 0)

		hp, _ := w.GetComponent(e, hpType)
		assertEqual(t, hp.(Health).HP, 5)
	})

	t.Run("update with no systems returns nil", func(t *testing.T) {
		t.Parallel()

		w := NewWorld()

		if err := w.Update(context.Background(), 0); err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})
}
