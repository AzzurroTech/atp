package storage

import (
	"sync"

	"github.com/AzzurroTech/atp/internal/models"
)

type Store interface {
	Add(source *models.Source)
	GetAll() []models.Source
	Exists(url string) bool
}

type MemoryStore struct {
	mu      sync.RWMutex
	sources map[string]*models.Source
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		sources: make(map[string]*models.Source),
	}
}

func (s *MemoryStore) Add(source *models.Source) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sources[source.ID] = source
}

func (s *MemoryStore) GetAll() []models.Source {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]models.Source, 0, len(s.sources))
	for _, v := range s.sources {
		result = append(result, *v)
	}
	return result
}

func (s *MemoryStore) Exists(url string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, v := range s.sources {
		if v.URL == url {
			return true
		}
	}
	return false
}
