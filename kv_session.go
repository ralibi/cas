package cas

import (
	"sync"
)

// KvSession
type KvSession struct {
	mu       sync.Mutex
	sessions map[string]string
}

func (s *KvSession) Read(key string) (string, bool) {
	return s.sessions[key], true
}

func (s *KvSession) Write(key, value string) {
	s.mu.Lock()
	s.sessions[key] = value
	s.mu.Unlock()
}

func (s *KvSession) Delete(key string) {
	s.mu.Lock()
	delete(s.sessions, key)
	s.mu.Unlock()
}
