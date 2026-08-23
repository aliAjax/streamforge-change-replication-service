package sourceinfra

import (
	"context"
	"fmt"
	domain "github.com/acme/streamforge-cdc/internal/source/domain"
	"sync"
)

type MemoryRepository struct {
	mu    sync.RWMutex
	items map[string]domain.Source
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{items: map[string]domain.Source{}}
}
func (r *MemoryRepository) Save(_ context.Context, s domain.Source) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items[s.ID] = s
	return nil
}
func (r *MemoryRepository) Get(_ context.Context, id string) (domain.Source, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, ok := r.items[id]
	if !ok {
		return s, fmt.Errorf("source %s not found", id)
	}
	return s, nil
}
func (r *MemoryRepository) List(_ context.Context) ([]domain.Source, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]domain.Source, 0, len(r.items))
	for _, s := range r.items {
		out = append(out, s)
	}
	return out, nil
}

type Validator struct{}

func (Validator) Validate(ctx context.Context, s domain.Source) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	if s.Endpoint == "" {
		return fmt.Errorf("endpoint is required")
	}
	return nil
}
