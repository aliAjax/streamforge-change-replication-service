package transactionapp

import (
	"context"
	"fmt"
	transaction "github.com/acme/streamforge-cdc/internal/transaction/domain"
)

type Reader interface {
	List(context.Context, string, int, int) ([]transaction.Event, error)
}
type ReplayService struct {
	Reader Reader
	Sink   Sink
}

func (s ReplayService) Replay(ctx context.Context, source string, cursor, limit int) (int, error) {
	if s.Reader == nil || s.Sink == nil {
		return 0, fmt.Errorf("replay dependencies are required")
	}
	events, err := s.Reader.List(ctx, source, cursor, limit)
	if err != nil {
		return 0, fmt.Errorf("replay read: %w", err)
	}
	for i, event := range events {
		tx := transaction.Transaction{ID: "replay-" + event.ID, Source: event.Source, Position: event.Position, Events: []transaction.Event{event}}
		if err := s.Sink.Persist(ctx, tx); err != nil {
			return i, fmt.Errorf("replay event %s: %w", event.ID, err)
		}
	}
	return len(events), nil
}
