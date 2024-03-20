package sync

import "sync"

type SafeMap[K comparable, V any] struct {
	data  map[K]V
	mutex sync.RWMutex
}

func (s *SafeMap[K, V]) Get(key K) (V, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	val, ok := s.data[key]
	return val, ok
}

func (s *SafeMap[K, V]) Put(key K, val V) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.data[key] = val
}

func (s *SafeMap[K, V]) LoadOrStore(key K, newVal V) (V, bool) {
	s.mutex.RLock()
	val, ok := s.data[key]
	if ok {
		return val, true
	}
	s.mutex.RUnlock()

	s.mutex.Lock()
	defer s.mutex.Unlock()

	// double check
	val, ok = s.data[key]
	if ok {
		return val, true
	}
	s.data[key] = newVal
	return newVal, false
}
