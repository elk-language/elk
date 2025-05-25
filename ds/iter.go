package ds

import "iter"

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
