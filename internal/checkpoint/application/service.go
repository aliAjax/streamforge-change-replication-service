package checkpointapp

import (
	"context"
	"fmt"
	domain "github.com/acme/streamforge-cdc/internal/checkpoint/domain"
	transaction "github.com/acme/streamforge-cdc/internal/transaction/domain"
)

type Service struct{ Store *domain.Store }

func (s *Service) Read(ctx context.Context, id string) (domain.State, error) {
	if s.Store == nil {
		return domain.State{}, fmt.Errorf("checkpoint store is nil")
	}
	return s.Store.Load(ctx, id)
}
func (s *Service) Mark(ctx context.Context, id, kind string, pos transaction.Position, version int64) error {
	return s.Store.Advance(ctx, id, kind, pos, version)
}
