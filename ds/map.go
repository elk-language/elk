package ds

// Map elements of one slice and create a new slice
// with those mapped values.
func MapSlice[T, R any](input []T, f func(T) R) []R {
	result := make([]R, len(input))
	for i, v := range input {
		result[i] = f(v)
	}
	return result
}
