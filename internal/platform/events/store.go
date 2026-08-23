package eventstore

import (
	"context"
	"fmt"
	transaction "github.com/acme/streamforge-cdc/internal/transaction/domain"
	"sync"
)

type Store struct {
	mu     sync.Mutex
	Txs    []transaction.Transaction
	Events []transaction.Event
}

func (s *Store) Persist(_ context.Context, t transaction.Transaction) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, old := range s.Txs {
		if old.ID == t.ID {
			return nil
		}
	}
	s.Txs = append(s.Txs, t)
	s.Events = append(s.Events, t.Events...)
	return nil
}
func (s *Store) List(after int, limit int) []transaction.Event {
	s.mu.Lock()
	defer s.mu.Unlock()
	if after < 0 {
		after = 0
	}
	if after >= len(s.Events) {
		return []transaction.Event{}
	}
	end := after + limit
	if limit <= 0 || end > len(s.Events) {
		end = len(s.Events)
	}
	return append([]transaction.Event(nil), s.Events[after:end]...)
}
func (s *Store) Count() int { s.mu.Lock(); defer s.mu.Unlock(); return len(s.Events) }
func (s *Store) Get(id string) (transaction.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, e := range s.Events {
		if e.ID == id {
			return e, nil
		}
	}
	return transaction.Event{}, fmt.Errorf("event %s not found", id)
}
