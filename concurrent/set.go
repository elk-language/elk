package concurrent

import (
	"iter"
	"maps"
	"sync"
)

type Set[V comparable] struct {
	Map map[V]struct{}
	mu  sync.RWMutex
}

func MakeSet[V comparable]() Set[V] {
	return Set[V]{
		Map: make(map[V]struct{}),
	}
}

func NewSet[V comparable]() *Set[V] {
	return &Set[V]{
		Map: make(map[V]struct{}),
	}
}

func NewSetWithValues[V comparable](m map[V]struct{}) *Set[V] {
	return &Set[V]{
		Map: m,
	}
}

func (s *Set[V]) All() iter.Seq[V] {
	return func(yield func(V) bool) {
		s.mu.Lock()
		defer s.mu.Unlock()

		for val := range s.Map {
			if !yield(val) {
				break
			}
		}
	}
}

func (s *Set[V]) Clone() *Set[V] {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.CloneUnsafe()
}

func (s *Set[V]) CloneValue() Set[V] {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.CloneUnsafeValue()
}

func (s *Set[V]) CloneUnsafe() *Set[V] {
	return &Set[V]{
		Map: maps.Clone(s.Map),
	}
}

func (s *Set[V]) CloneUnsafeValue() Set[V] {
	return Set[V]{
		Map: maps.Clone(s.Map),
	}
}

func (s *Set[V]) Len() int {
	return len(s.Map)
}

func (s *Set[V]) CopyElements(other *Set[V]) {
	s.mu.Lock()
	defer s.mu.Unlock()

	other.mu.Lock()
	defer other.mu.Unlock()

	maps.Copy(s.Map, other.Map)
}

func (s *Set[V]) CopyElementsUnsafe(other *Set[V]) {
	maps.Copy(s.Map, other.Map)
}

func (s *Set[V]) Remove(val V) {
	s.mu.Lock()
	delete(s.Map, val)
	s.mu.Unlock()
}

func (s *Set[V]) RemoveUnsafe(val V) {
	delete(s.Map, val)
}

func (s *Set[V]) Add(val V) {
	s.mu.Lock()
	s.Map[val] = struct{}{}
	s.mu.Unlock()
}

func (s *Set[V]) AddUnsafe(val V) {
	s.Map[val] = struct{}{}
}

func (s *Set[V]) Contains(val V) bool {
	s.mu.Lock()
	_, ok := s.Map[val]
	s.mu.Unlock()
	return ok
}

func (s *Set[V]) ContainsUnsafe(val V) bool {
	_, ok := s.Map[val]
	return ok
}
