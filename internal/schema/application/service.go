package schemaapp

import (
	"context"
	"fmt"
	domain "github.com/acme/streamforge-cdc/internal/schema/domain"
	"sync"
)

type Repository interface {
	Put(context.Context, domain.Table) error
	Latest(context.Context, string, string) (domain.Table, error)
}
type Service struct {
	repo     Repository
	mu       sync.Mutex
	strategy domain.Strategy
}

func New(r Repository, st domain.Strategy) *Service {
	if st == "" {
		st = domain.Compatible
	}
	return &Service{repo: r, strategy: st}
}
func (s *Service) Apply(ctx context.Context, t domain.Table) ([]string, error) {
	old, err := s.repo.Latest(ctx, t.Schema, t.Name)
	if err != nil {
		if err2 := s.repo.Put(ctx, t); err2 != nil {
			return nil, fmt.Errorf("store first schema: %w", err2)
		}
		return nil, nil
	}
	changes := domain.Compare(old, t)
	if len(changes) > 0 && s.strategy == domain.Strict {
		return changes, fmt.Errorf("incompatible schema change: %v", changes)
	}
	if err := s.repo.Put(ctx, t); err != nil {
		return changes, fmt.Errorf("store schema: %w", err)
	}
	return changes, nil
}
