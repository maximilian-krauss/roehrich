package utils

func Filter[T any](source []T, test func(T) bool) (destination []T) {
	for _, item := range source {
		if test(item) {
			destination = append(destination, item)
		}
	}
	return destination
}
