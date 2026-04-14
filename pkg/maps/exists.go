package maps

func Exists[T comparable, A any](m map[T]A, value T) bool {
	_, ok := m[value]
	return ok
}
