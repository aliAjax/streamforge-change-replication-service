package deliveryapp

import (
	"context"
	"fmt"
	buffer "github.com/acme/streamforge-cdc/internal/buffer/domain"
	delivery "github.com/acme/streamforge-cdc/internal/delivery/domain"
	"time"
)

type Service struct {
	Queue       *buffer.Queue
	Sender      delivery.Sender
	MaxAttempts int
	Batch       int
}

func (s *Service) Flush(ctx context.Context) (int, error) {
	if s.MaxAttempts < 1 {
		s.MaxAttempts = 3
	}
	if s.Batch < 1 {
		s.Batch = 10
	}
	sent := 0
	for i := 0; i < s.Batch; i++ {
		r, ok := s.Queue.Peek()
		if !ok {
			break
		}
		err := s.Sender.Send(ctx, []delivery.Message{{ID: r.ID, Transaction: r.Tx}})
		if err != nil {
			s.Queue.Fail(r.ID, s.MaxAttempts)
			return sent, fmt.Errorf("deliver %s: %w", r.ID, err)
		}
		s.Queue.Ack(r.ID)
		sent++
		time.Sleep(0)
	}
	return sent, nil
}
