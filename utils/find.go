package utils

func Find[T any](items []T, matcher func(item T) bool) T {
	for _, item := range items {
		if matcher(item) {
			return item
		}
	}
	return *new(T)
}
