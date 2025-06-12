package syncmap

import "sync"

type Map[K comparable, V any] struct {
	inner map[K]V
	mu    sync.RWMutex
}

func (m *Map[K, V]) Keys() []K {
	m.mu.RLock()
	defer m.mu.RUnlock()
	m.init()
	keys := make([]K, 0, len(m.inner))
	for k := range m.inner {
		keys = append(keys, k)
	}
	return keys
}

func (m *Map[K, V]) Get(k K) (V, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	m.init()
	v, ok := m.inner[k]
	return v, ok
}

func (m *Map[K, V]) Set(k K, v V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.init()
	m.inner[k] = v
}

func (m *Map[K, V]) Delete(k K) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.init()
	delete(m.inner, k)
}

func (m *Map[K, V]) Update(k K, u func(v V) V) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.init()
	v, ok := m.inner[k]
	if !ok {
		return false
	}
	m.inner[k] = u(v)
	return true
}

func (m *Map[K, V]) Transaction(fn func(m map[K]V) any) any {
	m.mu.Lock()
	defer m.mu.Unlock()
	return fn(m.inner)
}

func (m *Map[K, V]) init() {
	if m.inner == nil {
		m.inner = make(map[K]V)
	}
}
