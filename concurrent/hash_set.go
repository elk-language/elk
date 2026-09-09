package concurrent

import (
	"sync"

	"github.com/elk-language/elk/ds"
)

// Set with custom hashing and equality methods
type HashSet[V ds.Hashable] struct {
	Data ds.HashSet[V]
	mu   sync.RWMutex
}

func NewHashSet[V ds.Hashable]() *HashSet[V] {
	return &HashSet[V]{
		Data: ds.MakeHashSet[V](5),
	}
}

func NewHashSetWithCapacity[V ds.Hashable](cap int) *HashSet[V] {
	return &HashSet[V]{
		Data: ds.MakeHashSet[V](cap),
	}
}

func NewHashSetWithValues[V ds.Hashable](values []V) *HashSet[V] {
	return &HashSet[V]{
		Data: ds.MakeHashSetWithValues(values),
	}
}

func (s *HashSet[V]) Len() int {
	return s.Data.Len()
}

func (s *HashSet[V]) Capacity() int {
	return s.Data.Capacity()
}

func (s *HashSet[V]) LeftCapacity() int {
	return s.Data.LeftCapacity()
}

func (s *HashSet[V]) Remove(val V) bool {
	s.mu.Lock()
	result := s.Data.Remove(val)
	s.mu.Unlock()

	return result
}

func (s *HashSet[V]) RemoveUnsafe(val V) bool {
	return s.Data.Remove(val)
}

func (s *HashSet[V]) Add(val V) bool {
	s.mu.Lock()
	result := s.Data.Add(val)
	s.mu.Unlock()

	return result
}

func (s *HashSet[V]) AddUnsafe(val V) bool {
	return s.Data.Add(val)
}

func (s *HashSet[V]) SetCapacity(cap int) {
	s.mu.Lock()
	s.Data.SetCapacity(cap)
	s.mu.Unlock()
}

func (s *HashSet[V]) SetCapacityUnsafe(cap int) {
	s.Data.SetCapacity(cap)
}

func (s *HashSet[V]) Contains(val V) bool {
	s.mu.Lock()
	result := s.Data.Contains(val)
	s.mu.Unlock()

	return result
}

func (s *HashSet[V]) ContainsUnsafe(val V) bool {
	return s.Data.Contains(val)
}

func (s *HashSet[V]) Contains(val V) bool {
	s.mu.Lock()
	result := s.Data.Contains(val)
	s.mu.Unlock()

	return result
}

func (s *HashSet[V]) ContainsUnsafe(val V) bool {
	return s.Data.Contains(val)
}
