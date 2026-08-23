package bufferapp

import (
	"context"
	"fmt"
	buffer "github.com/acme/streamforge-cdc/internal/buffer/domain"
	transaction "github.com/acme/streamforge-cdc/internal/transaction/domain"
	"time"
)

type Service struct {
	Queue     *buffer.Queue
	Persisted func(context.Context, transaction.Position) error
	Confirmed func(context.Context, transaction.Position) error
}

func (s *Service) Append(ctx context.Context, tx transaction.Transaction) error {
	if err := s.Queue.Enqueue(tx); err != nil {
		return fmt.Errorf("append buffer: %w", err)
	}
	if s.Persisted != nil {
		if err := s.Persisted(ctx, tx.Position); err != nil {
			return fmt.Errorf("mark persisted: %w", err)
		}
	}
	return nil
}
func (s *Service) Drain(ctx context.Context, send func(context.Context, buffer.Record) error) (int, error) {
	count := 0
	for {
		rec, ok := s.Queue.Peek()
		if !ok {
			return count, nil
		}
		if rec.Dead {
			return count, fmt.Errorf("dead letter %s requires replay", rec.ID)
		}
		if err := send(ctx, rec); err != nil {
			s.Queue.Fail(rec.ID, 3)
			return count, err
		}
		s.Queue.Ack(rec.ID)
		if s.Confirmed != nil {
			if err := s.Confirmed(ctx, rec.Tx.Position); err != nil {
				return count, fmt.Errorf("mark confirmed: %w", err)
			}
		}
		count++
		select {
		case <-ctx.Done():
			return count, ctx.Err()
		default:
			time.Sleep(time.Millisecond)
		}
	}
}
