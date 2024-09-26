package concurrent

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
