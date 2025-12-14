package concurrent

import "sync"

type Map[K comparable, V any] struct {
	Map map[K]V
	mu  sync.RWMutex
}

func NewMap[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{
		Map: make(map[K]V),
	}
}

func NewMapWithValues[K comparable, V any](m map[K]V) *Map[K, V] {
	return &Map[K, V]{
		Map: m,
	}
}

func (m *Map[K, V]) Lock() {
	m.mu.Lock()
}

func (m *Map[K, V]) Unlock() {
	m.mu.Unlock()
}

func (m *Map[K, V]) Len() int {
	return len(m.Map)
}

func (m *Map[K, V]) Get(key K) (val V, ok bool) {
	m.mu.RLock()
	val, ok = m.Map[key]
	m.mu.RUnlock()
	return val, ok
}

func (m *Map[K, V]) GetUnsafe(key K) (val V, ok bool) {
	val, ok = m.Map[key]
	return val, ok
}

func (m *Map[K, V]) Delete(key K) {
	m.mu.Lock()
	delete(m.Map, key)
	m.mu.Unlock()
}

func (m *Map[K, V]) DeleteUnsafe(key K) {
	delete(m.Map, key)
}

func (m *Map[K, V]) Set(key K, value V) {
	m.mu.Lock()
	m.Map[key] = value
	m.mu.Unlock()
}

func (m *Map[K, V]) SetUnsafe(key K, value V) {
	m.Map[key] = value
}
