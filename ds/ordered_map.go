package ds

import "iter"

// OrderedMap implements a map with ordered keys.
// It iterates in insertion order.
type OrderedMap[K comparable, V any] struct {
	data  map[K]V
	order []K
}

// NewOrderedMap instantiates a new ordered map
func NewOrderedMap[K comparable, V any]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{
		data: make(map[K]V),
	}
}

// NewOrderedMap instantiates a new ordered map with the given capacity
func NewOrderedMapWithCap[K comparable, V any](cap int) *OrderedMap[K, V] {
	return &OrderedMap[K, V]{
		data:  make(map[K]V, cap),
		order: make([]K, 0, cap),
	}
}

// NewOrderedMap instantiates a new ordered map with the given pairs
func NewOrderedMapWithPairs[K comparable, V any](pairs ...Pair[K, V]) *OrderedMap[K, V] {
	m := NewOrderedMapWithCap[K, V](len(pairs))

	for _, pair := range pairs {
		m.Set(pair.Key, pair.Value)
	}

	return m
}

// Get return the value for the given key
func (m *OrderedMap[K, V]) Get(key K) V {
	return m.data[key]
}

// GetOK returns the value for the given key an a bool
// flag that is false when the key is not present
func (m *OrderedMap[K, V]) GetOk(key K) (V, bool) {
	val, ok := m.data[key]
	return val, ok
}

// Set adds a new key value pair or sets an existing
// key to the given value
func (m *OrderedMap[K, V]) Set(key K, val V) *OrderedMap[K, V] {
	_, present := m.data[key]
	m.data[key] = val
	if !present {
		m.order = append(m.order, key)
	}

	return m
}

// Includes checks if the given key exists within the map
func (m *OrderedMap[K, V]) Includes(key K) bool {
	_, present := m.data[key]
	return present
}

// Insert adds a new pair if the key is not already present in the map
// otherwise it does nothing and returns false
func (m *OrderedMap[K, V]) Insert(key K, val V) bool {
	_, present := m.data[key]
	if present {
		return false
	}

	m.Set(key, val)
	return true
}

// Delete removes the given key and returns true when the key was
// found
func (m *OrderedMap[K, V]) Delete(key K) bool {
	_, present := m.data[key]
	if !present {
		return false
	}

	delete(m.data, key)
	m.removeFromOrder(key)
	return true
}

func (m *OrderedMap[K, V]) removeFromOrder(key K) {
	for i, k := range m.order {
		if k == key {
			m.order = append(m.order[:i], m.order[i+1:]...)
		}
	}
}

// All returns an iterator that iterates over all
// key value pairs in insertion order
func (m *OrderedMap[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, key := range m.order {
			if !yield(key, m.data[key]) {
				return
			}
		}
	}
}
