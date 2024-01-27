package query

func SliceUnique[V comparable](s []V) []V {
	keys := make(map[V]bool, len(s))
	list := make([]V, 0, len(s))
	for _, entry := range s {
		if _, ok := keys[entry]; !ok {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func CopyMap[K comparable, V any](data map[K]V) map[K]V {
	newMap := make(map[K]V, len(data))

	for key, value := range data {
		newMap[key] = value
	}

	return newMap
}
