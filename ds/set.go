package ds

type Set[V comparable] map[V]struct{}

func MakeSet[V comparable](values ...V) Set[V] {
	set := make(Set[V], len(values))
	for _, value := range values {
		set.Add(value)
	}
	return set
}

func (s Set[V]) Add(val V) {
	s[val] = struct{}{}
}

func (s Set[V]) Remove(val V) {
	delete(s, val)
}

func (s Set[V]) Contains(val V) bool {
	_, ok := s[val]
	return ok
}
