package main

import "sync"

type SafeStorage struct {
	mu   sync.RWMutex
	data map[string]Record
}

func NewSafeMap() *SafeStorage {
	return &SafeStorage{
		data: make(map[string]Record),
	}
}

func (s *SafeStorage) Size() int {
	return len(s.data)
}

func (s *SafeStorage) Get(key string) (Record, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.data[key]
	return val, ok
}

func (s *SafeStorage) Set(key string, value Record) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

func (s *SafeStorage) GetKeys() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	keys := make([]string, 0, len(s.data))
	for key := range s.data {
		keys = append(keys, key)
	}
	return keys
}

func (s *SafeStorage) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
}
