package ds

type Pair[K comparable, V any] struct {
	Key   K
	Value V
}

func (p *Pair[K, V]) Copy() *Pair[K, V] {
	return &Pair[K, V]{
		Key:   p.Key,
		Value: p.Value,
	}
}

func NewPair[K comparable, V any](key K, val V) *Pair[K, V] {
	return &Pair[K, V]{
		Key:   key,
		Value: val,
	}
}

func MakePair[K comparable, V any](key K, val V) Pair[K, V] {
	return Pair[K, V]{
		Key:   key,
		Value: val,
	}
}
