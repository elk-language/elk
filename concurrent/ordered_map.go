package concurrent

import (
	"iter"
	"sync"

	"github.com/elk-language/elk/ds"
)

type OrderedMap[K comparable, V any] struct {
	Map ds.OrderedMap[K, V]
	mu  sync.RWMutex
}

func NewOrderedMap[K comparable, V any]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{
		Map: ds.MakeOrderedMap[K, V](),
	}
}

func NewOrderedMapWithPairs[K comparable, V any](pairs ...ds.Pair[K, V]) *OrderedMap[K, V] {
	m := ds.MakeOrderedMapWithCap[K, V](len(pairs))

	for _, pair := range pairs {
		m.Set(pair.Key, pair.Value)
	}

	return &OrderedMap[K, V]{
		Map: m,
	}
}

func (m *OrderedMap[K, V]) Lock() {
	m.mu.Lock()
}

func (m *OrderedMap[K, V]) Unlock() {
	m.mu.Unlock()
}

func (m *OrderedMap[K, V]) Len() int {
	return m.Map.Len()
}

func (m *OrderedMap[K, V]) Clear() {
	m.Lock()
	m.ClearUnsafe()
	m.Unlock()
}

func (m *OrderedMap[K, V]) ClearUnsafe() {
	m.Map = ds.MakeOrderedMap[K, V]()
}

func (m *OrderedMap[K, V]) Get(key K) (val V, ok bool) {
	m.mu.RLock()
	val, ok = m.Map.GetOk(key)
	m.mu.RUnlock()
	return val, ok
}

func (m *OrderedMap[K, V]) GetUnsafe(key K) (val V, ok bool) {
	val, ok = m.Map.GetOk(key)
	return val, ok
}

func (m *OrderedMap[K, V]) Delete(key K) bool {
	m.mu.Lock()
	result := m.Map.Delete(key)
	m.mu.Unlock()
	return result
}

func (m *OrderedMap[K, V]) DeleteUnsafe(key K) bool {
	return m.Map.Delete(key)
}

func (m *OrderedMap[K, V]) Set(key K, value V) {
	m.mu.Lock()
	m.Map.Set(key, value)
	m.mu.Unlock()
}

func (m *OrderedMap[K, V]) SetUnsafe(key K, value V) {
	m.Map.Set(key, value)
}

// All returns an iterator that iterates over all
// key value pairs in insertion order
func (m *OrderedMap[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		m.Lock()
		defer m.Unlock()

		for key, val := range m.All() {
			if !yield(key, val) {
				return
			}
		}
	}
}

// All returns an iterator that iterates over all
// keys in insertion order
func (m *OrderedMap[K, V]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		m.Lock()
		defer m.Unlock()

		for key := range m.Keys() {
			if !yield(key) {
				return
			}
		}
	}
}

// All returns an iterator that iterates over all
// values in insertion order
func (m *OrderedMap[K, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		m.Lock()
		defer m.Unlock()

		for val := range m.Values() {
			if !yield(val) {
				return
			}
		}
	}
}
