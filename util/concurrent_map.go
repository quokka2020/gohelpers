package util

import "sync"

type Map[K comparable, V any] interface {
	Size() int
	Get(K) (V, bool)
	Put(K, V) V
	ComputeIfAbsent(K, func() V) V
	Remove(K) V
	Clear()
	Empty() bool
}

type concurrent_map[K comparable, V any] struct {
	sync.Mutex
	data map[K]V
}

func NewConcurrentMap[K comparable, V any]() Map[K, V] {
	ret := concurrent_map[K, V]{
		data: map[K]V{},
	}
	return &ret
}

func (m *concurrent_map[K, V]) Size() int {
	m.Lock()
	defer m.Unlock()
	return len(m.data)
}

func (m *concurrent_map[K, V]) Empty() bool {
	return m.Size() == 0
}

func (m *concurrent_map[K, V]) Get(key K) (V, bool) {
	m.Lock()
	defer m.Unlock()
	current_value, found := m.data[key]
	return current_value, found
}

func (m *concurrent_map[K, V]) Put(key K, value V) V {
	m.Lock()
	defer m.Unlock()
	current_value := m.data[key]
	m.data[key] = value
	return current_value
}

func (m *concurrent_map[K, V]) ComputeIfAbsent(key K, f func() V) V {
	m.Lock()
	defer m.Unlock()
	current_value, found := m.data[key]
	if !found {
		new_value := f()
		m.data[key] = new_value
		return new_value
	}
	return current_value
}

func (m *concurrent_map[K, V]) Remove(key K) V {
	m.Lock()
	defer m.Unlock()
	current_value := m.data[key]
	delete(m.data, key)
	return current_value
}

func (m *concurrent_map[K, V]) Clear() {
	m.Lock()
	defer m.Unlock()
	m.data = map[K]V{}
}
