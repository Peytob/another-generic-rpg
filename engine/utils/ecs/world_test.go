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

const (
	posType  ComponentType = 1
	velType  ComponentType = 2
	hpType   ComponentType = 3
	manaType ComponentType = 4

	damageEventType EventType = 1
	healEventType   EventType = 2
)

func (Position) Type() ComponentType { return posType }
func (Velocity) Type() ComponentType { return velType }
func (Health) Type() ComponentType   { return hpType }
func (Mana) Type() ComponentType     { return manaType }

func (*DamageEvent) Type() EventType { return damageEventType }
func (*HealEvent) Type() EventType   { return healEventType }

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
