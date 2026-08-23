package checkpointdomain

import (
	"context"
	"fmt"
	transaction "github.com/acme/streamforge-cdc/internal/transaction/domain"
	"sync"
)

type State struct {
	Source    string               `json:"source"`
	Read      transaction.Position `json:"read"`
	Persisted transaction.Position `json:"persisted"`
	Confirmed transaction.Position `json:"confirmed"`
	Version   int64                `json:"version"`
}
type Store struct {
	mu    sync.Mutex
	items map[string]State
}

func New() *Store { return &Store{items: map[string]State{}} }
func (s *Store) Load(_ context.Context, id string) (State, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.items[id]
	if !ok {
		return State{Source: id}, nil
	}
	return v, nil
}
func (s *Store) Advance(_ context.Context, id string, kind string, pos transaction.Position, version int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	v := s.items[id]
	if version != 0 && v.Version != version {
		return fmt.Errorf("checkpoint %s version conflict", id)
	}
	switch kind {
	case "read":
		v.Read = pos
	case "persisted":
		v.Persisted = pos
	case "confirmed":
		v.Confirmed = pos
	default:
		return fmt.Errorf("unknown checkpoint kind %s", kind)
	}
	v.Source = id
	v.Version++
	s.items[id] = v
	return nil
}
func (s *Store) Snapshot(_ context.Context) []State {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]State, 0, len(s.items))
	for _, v := range s.items {
		out = append(out, v)
	}
	return out
}
