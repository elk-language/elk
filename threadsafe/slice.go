package threadsafe

import "sync"

type Slice[V any] struct {
	Slice []V
	mu    sync.RWMutex
}

func NewSlice[V any]() *Slice[V] {
	return &Slice[V]{}
}

func NewSliceWithValues[V any](s []V) *Slice[V] {
	return &Slice[V]{
		Slice: s,
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

func (s *Slice[V]) Set(index int, val V) (ok bool) {
	s.mu.Lock()
	if index < len(s.Slice) && index > 0 {
		s.Slice[index] = val
		ok = true
	}
	s.mu.Unlock()
	return ok
}

func (s *Slice[V]) Append(val V) {
	s.mu.Lock()
	s.Slice = append(s.Slice, val)
	s.mu.Unlock()
}
