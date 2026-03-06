package scheduling

import "sync"

type PreviewStore[T any] interface {
	Save(id string, payload T)
	Get(id string) (T, bool)
	GetLatest() (T, bool)
}

type inMemoryPreviewStore[T any] struct {
	mu       sync.RWMutex
	items    map[string]T
	latestID string
}

func NewPreviewStore[T any]() PreviewStore[T] {
	return &inMemoryPreviewStore[T]{
		items: make(map[string]T),
	}
}

func (s *inMemoryPreviewStore[T]) Save(id string, payload T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.items[id] = payload
	s.latestID = id
}

func (s *inMemoryPreviewStore[T]) Get(id string) (T, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	payload, ok := s.items[id]
	return payload, ok
}

func (s *inMemoryPreviewStore[T]) GetLatest() (T, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	payload, ok := s.items[s.latestID]
	return payload, ok
}
