package kit

// Map transforms a slice into a new slice of the return type of the function
// parameter.
func Map[E, T any](s []E, f func(E) T) []T {
	result := make([]T, 0, len(s))

	for _, val := range s {
		result = append(result, f(val))
	}
	return result
}
