package kit

func Abs(i int) int {
	if i < 0 {
		i *= -1
	}
	return i
}

type Summable interface {
	int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 |
		float32 | float64
}

func Sum[T Summable](nums []T) T {
	var out T
	for _, n := range nums {
		out += n
	}
	return out
}

func Mult[T Summable](nums []T) T {
	out := nums[0]
	for _, n := range nums[1:] {
		out *= n
	}
	return out
}
