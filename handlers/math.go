package handlers

type number interface {
	int | float64
}

func min[T number](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func max[T number](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func tern[T any](cond bool, a, b T) T {
	if cond {
		return a
	}
	return b
}
