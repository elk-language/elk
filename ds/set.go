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

// Returns a new set containing the elements of all given sets
func (s Set[V]) Concat(sets ...Set[V]) Set[V] {
	newSet := make(Set[V])

	for element := range s {
		newSet.Add(element)
	}

	for _, set := range sets {
		for element := range set {
			newSet.Add(element)
		}
	}

	return newSet
}

// Mutate the set by adding the elements of all given sets
func (s Set[V]) ConcatMut(sets ...Set[V]) {
	for _, set := range sets {
		for element := range set {
			s.Add(element)
		}
	}
}
