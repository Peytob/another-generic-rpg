package ecs

import (
	"sort"
	"testing"
)

type Position struct{ X, Y int }
type Velocity struct{ DX, DY int }
type Health struct{ HP int }
type Mana struct{ MP int }

type DamageEvent struct{ Amount int }
type HealEvent struct{ Amount int }

var (
	posType  = ComponentTypeOf(Position{})
	velType  = ComponentTypeOf(Velocity{})
	hpType   = ComponentTypeOf(Health{})
	manaType = ComponentTypeOf(Mana{})
)

func assertEqual[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func assertSliceEqual[T comparable](t *testing.T, got, want []T) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("len = %d, want %d (got %v, want %v)", len(got), len(want), got, want)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("index %d: got %v, want %v", i, got[i], want[i])
		}
	}
}

func sortedEntities(entities []Entity) []Entity {
	sorted := make([]Entity, len(entities))
	copy(sorted, entities)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })
	return sorted
}
