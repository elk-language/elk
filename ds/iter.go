package ds

import (
	"iter"
)

// Map elements of one slice and create a new slice
// with those mapped values.
func MapSlice[T, R any](input []T, f func(T) R) []R {
	result := make([]R, len(input))
	for i, v := range input {
		result[i] = f(v)
	}
	return result
}

// Iterate over the elements of a slice in reverse order.
func ReverseSlice[V any](s []V) iter.Seq[V] {
	return func(yield func(V) bool) {
		for i := range len(s) {
			v := s[len(s)-i-1]
			if !yield(v) {
				return
			}
		}
	}
}

// Reverse a Seq2
func ReverseSeq2[K comparable, V any](i iter.Seq2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		var buff []Pair[K, V]

		for k, v := range i {
			buff = append(buff, MakePair(k, v))
		}

		for i := len(buff) - 1; i >= 0; i-- {
			if !yield(buff[i].Key, buff[i].Value) {
				return
			}
		}
	}
}

// Reverse a Seq
func ReverseSeq[V any](i iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		var buff []V

		for v := range i {
			buff = append(buff, v)
		}

		for i := len(buff) - 1; i >= 0; i-- {
			if !yield(buff[i]) {
				return
			}
		}
	}
}
