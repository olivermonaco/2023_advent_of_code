package kit

func Ptr[T any](o T) *T {
	return &o
}

func Deref[T any](ptr *T) T {
	if ptr != nil {
		return *ptr
	}

	var zeroValue T
	return zeroValue
}
