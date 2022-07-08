package syncmap

import (
	"sync"
)

type Map[K comparable, V any] struct {
	items map[K]V
	sync.Mutex
}

func New[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{
		items: make(map[K]V),
	}
}

func (m *Map[K, V]) Store(key K, value V) {
	m.Lock()
	m.items[key] = value
	m.Unlock()
}

func (m *Map[K, V]) Load(key K) (V, bool) {
	m.Lock()
	val, ok := m.items[key]
	m.Unlock()

	return val, ok
}

func (m *Map[K, V]) Items() map[K]V {
	items := make(map[K]V)

	m.Lock()
	for key, item := range m.items {
		items[key] = item
	}
	m.Unlock()

	return items
}
