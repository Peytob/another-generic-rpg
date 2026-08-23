package hid

import "fmt"

type bindings[K comparable, C any] struct {
	name     string
	bindings map[K]C
}

func newBindings[K comparable, C any](name string) bindings[K, C] {
	return bindings[K, C]{
		name:     name,
		bindings: make(map[K]C),
	}
}

func (b bindings[K, C]) add(key K, callback C) error {
	if _, exists := b.bindings[key]; exists {
		return fmt.Errorf("binding already exists: %v", key)
	}
	b.bindings[key] = callback
	return nil
}

func (b bindings[K, C]) lookup(key K) (C, bool) {
	callback, ok := b.bindings[key]
	return callback, ok
}

func (b bindings[K, C]) keys(yield func(K) bool) {
	for key := range b.bindings {
		if !yield(key) {
			return
		}
	}
}

func (b bindings[K, C]) Name() string {
	return b.name
}
