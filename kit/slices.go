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

func Reduce[E, T any](s []E, initializer T, f func(T, E) T) T {
	r := initializer
	for _, val := range s {
		r = f(r, val)
	}
	return r
}

func Filter[T any](s []T, f func(T) bool) []T {
	nS := s[:0]
	for _, x := range s {
		if f(x) {
			nS = append(nS, x)
		}
	}
	return nS
}

func Keys[C comparable, T any](m map[C]T) []C {
	keys := make([]C, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
