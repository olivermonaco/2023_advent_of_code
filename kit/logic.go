package kit

func IsBetween(n, a, b int) bool {
	min := a
	max := b
	if a > b {
		min = b
		max = a
	}
	return min <= n && n <= max
}
