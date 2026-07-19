package ecs

import (
	"testing"
)

func TestRegisterComponent(t *testing.T) {
	t.Parallel()

	t.Run("adds component to entity", func(t *testing.T) {
		t.Parallel()

		w := NewWorld()
		e := w.NewEntity()
		w.RegisterComponent(e, Position{10, 20})

		c, ok := w.GetComponent(e, posType)
		if !ok {
			t.Fatal("expected component to be found")
		}
		assertEqual(t, c.(Position), Position{10, 20})
	})

	t.Run("replaces existing component of same type", func(t *testing.T) {
		t.Parallel()

		w := NewWorld()
		e := w.NewEntity()
		w.RegisterComponent(e, Position{1, 1})
		w.RegisterComponent(e, Position{2, 2})

		c, _ := w.GetComponent(e, posType)
		assertEqual(t, c.(Position), Position{2, 2})

		got := w.Query(posType)
		assertSliceEqual(t, got, []Entity{e})
	})

	t.Run("allows multiple component types per entity", func(t *testing.T) {
		t.Parallel()

		w := NewWorld()
		e := w.NewEntity()
		w.RegisterComponent(e, Position{0, 0})
		w.RegisterComponent(e, Velocity{1, 1})
		w.RegisterComponent(e, Health{50})

		if !w.HasComponent(e, posType) {
			t.Error("missing Position")
		}
		if !w.HasComponent(e, velType) {
			t.Error("missing Velocity")
		}
		if !w.HasComponent(e, hpType) {
			t.Error("missing Health")
		}
	})

	t.Run("on non-existent entity is a no-op", func(t *testing.T) {
		t.Parallel()

		w := NewWorld()

		w.RegisterComponent(999, Position{0, 0})

		got := w.Query(posType)
		if len(got) != 0 {
			t.Errorf("expected no results, got %v", got)
		}
	})

	t.Run("on removed entity is a no-op", func(t *testing.T) {
		t.Parallel()

		w := NewWorld()
		e := w.NewEntity()
		w.RemoveEntity(e)

		w.RegisterComponent(e, Position{0, 0})

		if w.HasComponent(e, posType) {
			t.Error("component was registered for removed entity")
		}
		got := w.Query(posType)
		if len(got) != 0 {
			t.Errorf("expected no query results, got %v", got)
		}
	})

	t.Run("on InvalidEntity is a no-op", func(t *testing.T) {
		t.Parallel()

		w := NewWorld()
		w.RegisterComponent(InvalidEntity, Position{0, 0})

		if w.HasComponent(InvalidEntity, posType) {
			t.Error("component was registered for InvalidEntity")
		}
	})
}

func TestUnregisterComponent(t *testing.T) {
	t.Parallel()

	t.Run("removes component from entity", func(t *testing.T) {
		t.Parallel()

		w := NewWorld()
		e := w.NewEntity()
		w.RegisterComponent(e, Position{0, 0})
		w.RegisterComponent(e, Health{10})

		w.UnregisterComponent(e, posType)

		if w.HasComponent(e, posType) {
			t.Error("Position was not removed")
		}
		if !w.HasComponent(e, hpType) {
			t.Error("Health should still be present")
		}
	})

	t.Run("entity no longer appears in query", func(t *testing.T) {
		t.Parallel()

		w := NewWorld()
		e1 := w.NewEntity()
		e2 := w.NewEntity()
		w.RegisterComponent(e1, Position{0, 0})
		w.RegisterComponent(e2, Position{1, 1})

		w.UnregisterComponent(e1, posType)

		got := w.Query(posType)
		assertSliceEqual(t, got, []Entity{e2})
	})

	t.Run("unregistering non-present type is safe", func(t *testing.T) {
		t.Parallel()

		w := NewWorld()
		e := w.NewEntity()
		w.RegisterComponent(e, Position{0, 0})

		w.UnregisterComponent(e, velType)

		if !w.HasComponent(e, posType) {
			t.Error("Position should still be present")
		}
	})

	t.Run("on non-existent entity is a no-op", func(t *testing.T) {
		t.Parallel()

		w := NewWorld()

		w.UnregisterComponent(999, posType)

		got := w.Query(posType)
		if len(got) != 0 {
			t.Errorf("expected no results, got %v", got)
		}
	})
}

func TestGetComponent(t *testing.T) {
	t.Parallel()

	t.Run("returns component and true", func(t *testing.T) {
		t.Parallel()

		w := NewWorld()
		e := w.NewEntity()
		w.RegisterComponent(e, Position{5, 10})

		c, ok := w.GetComponent(e, posType)
		if !ok {
			t.Fatal("expected ok=true")
		}
		assertEqual(t, c.(Position), Position{5, 10})
	})

	t.Run("returns nil and false for missing component", func(t *testing.T) {
		t.Parallel()

		w := NewWorld()
		e := w.NewEntity()

		c, ok := w.GetComponent(e, posType)
		if ok {
			t.Error("expected ok=false")
		}
		if c != nil {
			t.Errorf("expected nil, got %v", c)
		}
	})

	t.Run("returns nil and false for non-existent entity", func(t *testing.T) {
		t.Parallel()

		w := NewWorld()

		c, ok := w.GetComponent(999, posType)
		if ok {
			t.Error("expected ok=false")
		}
		if c != nil {
			t.Errorf("expected nil, got %v", c)
		}
	})
}

func TestHasComponent(t *testing.T) {
	t.Parallel()

	t.Run("true when present", func(t *testing.T) {
		t.Parallel()

		w := NewWorld()
		e := w.NewEntity()
		w.RegisterComponent(e, Position{0, 0})

		if !w.HasComponent(e, posType) {
			t.Error("expected HasComponent=true")
		}
	})

	t.Run("false when absent", func(t *testing.T) {
		t.Parallel()

		w := NewWorld()
		e := w.NewEntity()

		if w.HasComponent(e, posType) {
			t.Error("expected HasComponent=false")
		}
	})

	t.Run("false for non-existent entity", func(t *testing.T) {
		t.Parallel()

		w := NewWorld()

		if w.HasComponent(999, posType) {
			t.Error("expected HasComponent=false for non-existent entity")
		}
	})
}

func TestQuery(t *testing.T) {
	t.Parallel()

	setupWorld := func() (World, Entity, Entity, Entity) {
		w := NewWorld()
		e1 := w.NewEntity()
		e2 := w.NewEntity()
		e3 := w.NewEntity()

		w.RegisterComponent(e1, Position{0, 0})
		w.RegisterComponent(e1, Velocity{1, 0})

		w.RegisterComponent(e2, Position{1, 1})
		w.RegisterComponent(e2, Velocity{0, 1})
		w.RegisterComponent(e2, Health{50})

		// e3 has only Health
		w.RegisterComponent(e3, Health{100})

		return w, e1, e2, e3
	}

	t.Run("single component type", func(t *testing.T) {
		t.Parallel()

		w, e1, e2, _ := setupWorld()

		got := sortedEntities(w.Query(posType))
		assertSliceEqual(t, got, []Entity{e1, e2})
	})

	t.Run("multiple component types (intersection)", func(t *testing.T) {
		t.Parallel()

		w, e1, e2, e3 := setupWorld()

		got := sortedEntities(w.Query(posType, velType))
		assertSliceEqual(t, got, []Entity{e1, e2})

		got = sortedEntities(w.Query(posType, velType, hpType))
		assertSliceEqual(t, got, []Entity{e2})

		got = sortedEntities(w.Query(hpType))
		assertSliceEqual(t, got, []Entity{e2, e3})
	})

	t.Run("no matches returns empty slice", func(t *testing.T) {
		t.Parallel()

		w, _, _, _ := setupWorld()

		got := w.Query(posType, manaType)
		if len(got) != 0 {
			t.Errorf("expected no matches, got %v", got)
		}
	})

	t.Run("no component types returns nil", func(t *testing.T) {
		t.Parallel()

		w, _, _, _ := setupWorld()

		if got := w.Query(); got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})

	t.Run("excludes entity removed after registration", func(t *testing.T) {
		t.Parallel()

		w, e1, e2, _ := setupWorld()

		w.RemoveEntity(e1)

		got := sortedEntities(w.Query(posType))
		assertSliceEqual(t, got, []Entity{e2})
	})
}

func TestRegisterUnregisterQueryConsistency(t *testing.T) {
	t.Parallel()

	t.Run("unregister then re-register works correctly", func(t *testing.T) {
		t.Parallel()

		w := NewWorld()
		e := w.NewEntity()

		w.RegisterComponent(e, Position{1, 1})
		w.UnregisterComponent(e, posType)
		w.RegisterComponent(e, Position{2, 2})

		c, ok := w.GetComponent(e, posType)
		if !ok {
			t.Fatal("expected component after re-register")
		}
		assertEqual(t, c.(Position), Position{2, 2})

		got := w.Query(posType)
		assertSliceEqual(t, got, []Entity{e})
	})
}
