package store

import "sync"

type Store interface {
	SetID(name, id string)
	GetID(name string) *string
}

type InMemoryStore struct {
	sync.Mutex
	ids map[string]string
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		ids: make(map[string]string),
	}
}

func (s *InMemoryStore) SetID(name, id string) {
	s.Lock()
	defer s.Unlock()

	s.ids[name] = id
}

func (s *InMemoryStore) GetID(name string) *string {
	s.Lock()
	defer s.Unlock()

	id, exists := s.ids[name]
	if !exists {
		return nil
	}
	return &id
}
