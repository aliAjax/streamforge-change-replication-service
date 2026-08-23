package transactionapp

import (
	"context"
	"fmt"
	"sync"
	"time"

	domain "github.com/acme/streamforge-cdc/internal/transaction/domain"
)

type Sink interface {
	Persist(context.Context, domain.Transaction) error
}

type Service struct {
	mu     sync.Mutex
	active map[string]*domain.Transaction
	sink   Sink
	clock  func() time.Time
}

func New(s Sink) *Service {
	return &Service{active: make(map[string]*domain.Transaction), sink: s, clock: time.Now}
}

func (s *Service) Begin(ctx context.Context, source, id string, pos domain.Position) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.active[id]; ok {
		return fmt.Errorf("begin %s: already active", id)
	}
	s.active[id] = &domain.Transaction{ID: id, Source: source, Position: pos, StartedAt: s.clock()}
	return nil
}

func (s *Service) Append(ctx context.Context, id string, e domain.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.active[id]
	if !ok {
		return fmt.Errorf("append %s: transaction not found", id)
	}
	if err := t.Append(e); err != nil {
		return fmt.Errorf("append %s: %w", id, err)
	}
	return nil
}

func (s *Service) Commit(ctx context.Context, id string) (domain.Transaction, error) {
	s.mu.Lock()
	t, ok := s.active[id]
	if !ok {
		s.mu.Unlock()
		return domain.Transaction{}, fmt.Errorf("commit %s: transaction not found", id)
	}
	if err := t.Commit(s.clock()); err != nil {
		s.mu.Unlock()
		return domain.Transaction{}, fmt.Errorf("commit %s: %w", id, err)
	}
	delete(s.active, id)
	copyTx := *t
	copyTx.Events = append([]domain.Event(nil), t.Events...)
	s.mu.Unlock()
	if s.sink != nil {
		if err := s.sink.Persist(ctx, copyTx); err != nil {
			return domain.Transaction{}, fmt.Errorf("persist %s: %w", id, err)
		}
	}
	return copyTx, nil
}

func (s *Service) Abort(_ context.Context, id string) {
	s.mu.Lock()
	delete(s.active, id)
	s.mu.Unlock()
}
func (s *Service) Active() int { s.mu.Lock(); defer s.mu.Unlock(); return len(s.active) }
