package repository

import (
	"maps"
)

type Index[K comparable, V any] struct {
	index map[K]V
	all   []V
}

func NewIndex[K comparable, V any]() Index[K, V] {
	return Index[K, V]{
		index: make(map[K]V),
	}
}

func (i Index[K, V]) Add(key K, value V) bool {
	_, exists := i.index[key]
	if exists {
		return false
	}

	i.index[key] = value
	i.all = append(i.all, value)

	return true
}

func (i Index[K, V]) Remove(key K) bool {
	_, exists := i.index[key]
	if exists {
		return false
	}

	delete(i.index, key)
	i.all = nil

	return true
}

func (i Index[K, V]) Get(key K) (V, bool) {
	v, ok := i.index[key]
	return v, ok
}

func (i Index[K, V]) Array() []V {
	if i.all == nil {
		i.updateAllCache()
	}

	return i.all
}

func (i Index[K, V]) Count() int {
	return len(i.all)
}

func (i Index[K, V]) Clear() {
	clear(i.index)
	i.all = make([]V, 0)
}

func (i Index[K, V]) updateAllCache() {
	i.all = make([]V, 0, len(i.index))

	for v := range maps.Values(i.index) {
		i.all = append(i.all, v)
	}
}
