package query

import "sync"

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

//nolint:exhaustruct
type syncmap[K comparable, V any] struct {
	m sync.Map
}

func (m *syncmap[K, V]) Delete(key K) { m.m.Delete(key) }
func (m *syncmap[K, V]) Load(key K) (value V, ok bool) {
	v, ok := m.m.Load(key)
	if !ok {
		return value, ok
	}
	return v.(V), ok
}
func (m *syncmap[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	v, loaded := m.m.LoadAndDelete(key)
	if !loaded {
		return value, loaded
	}
	return v.(V), loaded
}
func (m *syncmap[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	a, loaded := m.m.LoadOrStore(key, value)
	return a.(V), loaded
}
func (m *syncmap[K, V]) Range(f func(key K, value V) bool) {
	m.m.Range(func(key, value any) bool { return f(key.(K), value.(V)) })
}
func (m *syncmap[K, V]) Store(key K, value V) { m.m.Store(key, value) }
