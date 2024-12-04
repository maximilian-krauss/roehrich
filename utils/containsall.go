package utils

func ContainsAll[T comparable](subset, superset []T) bool {
	elementMap := make(map[T]bool)
	for _, elem := range superset {
		elementMap[elem] = true
	}

	for _, elem := range subset {
		if !elementMap[elem] {
			return false
		}
	}
	return true
}
