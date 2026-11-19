package concurrent

import (
	"iter"
	"slices"
	"sync"
)

type Slice[V any] struct {
	Slice []V
	mu    sync.RWMutex
}

func MakeSlice[V any]() Slice[V] {
	return Slice[V]{}
}

func NewSlice[V any]() *Slice[V] {
	return &Slice[V]{}
}

func NewSliceWithValues[V any](s []V) *Slice[V] {
	return &Slice[V]{
		Slice: s,
	}
}

func (s *Slice[V]) All() iter.Seq2[int, V] {
	return func(yield func(int, V) bool) {
		s.mu.Lock()
		defer s.mu.Unlock()

		for i, val := range s.Slice {
			if !yield(i, val) {
				break
			}
		}
	}
}

func (s *Slice[V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		s.mu.Lock()
		defer s.mu.Unlock()

		for _, val := range s.Slice {
			if !yield(val) {
				break
			}
		}
	}
}

func (s *Slice[V]) Clone() *Slice[V] {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.CloneUnsafe()
}

func (s *Slice[V]) CloneValue() Slice[V] {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.CloneUnsafeValue()
}

func (s *Slice[V]) CloneUnsafe() *Slice[V] {
	return &Slice[V]{
		Slice: slices.Clone(s.Slice),
	}
}

func (s *Slice[V]) CloneUnsafeValue() Slice[V] {
	return Slice[V]{
		Slice: slices.Clone(s.Slice),
	}
}

func (s *Slice[V]) Len() int {
	return len(s.Slice)
}

func (s *Slice[V]) Cap() int {
	return cap(s.Slice)
}

func (s *Slice[V]) Get(index int) (val V, ok bool) {
	s.mu.RLock()
	if index < len(s.Slice) && index > 0 {
		val = s.Slice[index]
		ok = true
	}
	s.mu.RUnlock()
	return val, ok
}

func (s *Slice[V]) GetUnsafe(index int) (val V, ok bool) {
	if index < len(s.Slice) && index > 0 {
		val = s.Slice[index]
		ok = true
	}
	return val, ok
}

func (s *Slice[V]) Set(index int, val V) (ok bool) {
	s.mu.Lock()
	if index < len(s.Slice) && index > 0 {
		s.Slice[index] = val
		ok = true
	}
	s.mu.Unlock()
	return ok
}

func (s *Slice[V]) SetUnsafe(index int, val V) (ok bool) {
	if index < len(s.Slice) && index > 0 {
		s.Slice[index] = val
		ok = true
	}
	return ok
}

func (s *Slice[V]) Push(val V) {
	s.mu.Lock()
	s.Slice = append(s.Slice, val)
	s.mu.Unlock()
}

func (s *Slice[V]) PushUnsafe(val V) {
	s.Slice = append(s.Slice, val)
}

func (s *Slice[V]) Append(values ...V) {
	s.mu.Lock()
	s.Slice = append(s.Slice, values...)
	s.mu.Unlock()
}

func (s *Slice[V]) AppendUnsafe(values ...V) {
	s.Slice = append(s.Slice, values...)
}
