package ecs

import (
	"testing"
)

func TestNewEntity(t *testing.T) {
	t.Parallel()

	t.Run("first entity is not InvalidEntity", func(t *testing.T) {
		t.Parallel()

		w := NewWorld()
		e := w.NewEntity()

		assertEqual(t, e, Entity(1))
		if e == InvalidEntity {
			t.Error("first entity should not be InvalidEntity (0)")
		}
	})

	t.Run("each entity is unique and monotonically increasing", func(t *testing.T) {
		t.Parallel()

		w := NewWorld()

		var ids []Entity
		for i := 0; i < 100; i++ {
			ids = append(ids, w.NewEntity())
		}

		for i, id := range ids {
			want := Entity(i + 1)
			if id != want {
				t.Errorf("entity %d = %d, want %d", i, id, want)
			}
		}
	})
}

func TestEntities(t *testing.T) {
	t.Parallel()

	t.Run("returns all active entities", func(t *testing.T) {
		t.Parallel()

		w := NewWorld()
		e1 := w.NewEntity()
		e2 := w.NewEntity()
		e3 := w.NewEntity()

		got := sortedEntities(w.Entities())

		assertSliceEqual(t, got, []Entity{e1, e2, e3})
	})

	t.Run("empty world returns empty slice", func(t *testing.T) {
		t.Parallel()

		w := NewWorld()

		if entities := w.Entities(); len(entities) != 0 {
			t.Errorf("expected empty slice, got %v", entities)
		}
	})

	t.Run("returned slice is a copy", func(t *testing.T) {
		t.Parallel()

		w := NewWorld()
		w.NewEntity()
		w.NewEntity()

		entities := w.Entities()
		entities[0] = 9999

		again := w.Entities()
		if again[0] == Entity(9999) {
			t.Error("Entities() should return a copy")
		}
	})
}

func TestRemoveEntity(t *testing.T) {
	t.Parallel()

	t.Run("removes entity from Entities list", func(t *testing.T) {
		t.Parallel()

		w := NewWorld()
		e1 := w.NewEntity()
		e2 := w.NewEntity()

		w.RemoveEntity(e1)

		got := w.Entities()
		assertSliceEqual(t, got, []Entity{e2})
	})

	t.Run("removes all components from entity", func(t *testing.T) {
		t.Parallel()

		w := NewWorld()
		e := w.NewEntity()
		w.RegisterComponent(e, Position{1, 2})
		w.RegisterComponent(e, Health{100})

		w.RemoveEntity(e)

		if w.HasComponent(e, posType) {
			t.Error("entity still has Position after removal")
		}
		if w.HasComponent(e, hpType) {
			t.Error("entity still has Health after removal")
		}
	})

	t.Run("removes entity from query index", func(t *testing.T) {
		t.Parallel()

		w := NewWorld()
		e1 := w.NewEntity()
		e2 := w.NewEntity()
		w.RegisterComponent(e1, Position{0, 0})
		w.RegisterComponent(e2, Position{1, 1})

		w.RemoveEntity(e1)

		got := w.Query(posType)
		assertSliceEqual(t, got, []Entity{e2})
	})

	t.Run("remove non-existent entity is safe", func(t *testing.T) {
		t.Parallel()

		w := NewWorld()
		e := w.NewEntity()

		w.RemoveEntity(e + 999)

		if len(w.Entities()) != 1 {
			t.Error("removing non-existent entity affected state")
		}
	})

	t.Run("double remove is safe", func(t *testing.T) {
		t.Parallel()

		w := NewWorld()
		e := w.NewEntity()
		w.RegisterComponent(e, Position{0, 0})

		w.RemoveEntity(e)
		w.RemoveEntity(e)

		if len(w.Entities()) != 0 {
			t.Error("double remove left entities behind")
		}
	})
}
