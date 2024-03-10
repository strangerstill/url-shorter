package app

import "sync"

type Storage struct {
	urls map[string]string
	mu   sync.RWMutex
}

func NewStorage() *Storage {
	return &Storage{
		urls: make(map[string]string),
	}
}

func (s *Storage) Set(id, url string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.urls[id]
	if ok {
		return false
	}
	s.urls[id] = url
	return true
}

func (s *Storage) Get(id string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	url, ok := s.urls[id]
	return url, ok
}
