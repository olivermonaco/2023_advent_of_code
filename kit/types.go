package kit

func Ptr[T any](o T) *T {
	return &o
}
