package schemainfra

import (
	"context"
	"fmt"
	domain "github.com/acme/streamforge-cdc/internal/schema/domain"
	"sync"
)

type MemoryRepository struct {
	mu    sync.RWMutex
	items map[string]domain.Table
}

func New() *MemoryRepository { return &MemoryRepository{items: map[string]domain.Table{}} }
func (r *MemoryRepository) Put(_ context.Context, t domain.Table) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if t.Version == 0 {
		t.Version = 1
	}
	r.items[t.Schema+"."+t.Name] = t
	return nil
}
func (r *MemoryRepository) Latest(_ context.Context, s, n string) (domain.Table, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.items[s+"."+n]
	if !ok {
		return t, fmt.Errorf("schema %s.%s not found", s, n)
	}
	return t, nil
}
