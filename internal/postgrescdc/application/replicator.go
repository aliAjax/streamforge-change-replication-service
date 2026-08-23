package postgresapp

import (
	"context"
	"fmt"
	postgres "github.com/acme/streamforge-cdc/internal/postgrescdc/domain"
	transaction "github.com/acme/streamforge-cdc/internal/transaction/domain"
	"io"
	"time"
)

type Stream interface {
	Read(context.Context) (io.Reader, error)
	Confirm(context.Context, transaction.Position) error
}
type Decoder struct{ MaxMessage int }

func (d Decoder) Run(ctx context.Context, stream Stream, emit func(transaction.Event) error) error {
	r, err := stream.Read(ctx)
	if err != nil {
		return fmt.Errorf("postgres stream: %w", err)
	}
	if d.MaxMessage == 0 {
		d.MaxMessage = 16 << 20
	}
	for {
		m, err := postgres.Decode(r, d.MaxMessage)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("decode postgres: %w", err)
		}
		if m.Tag == 'k' {
			continue
		}
		if m.Tag == 'c' {
			e := transaction.Event{Source: "postgres", Operation: transaction.Insert, TransactionID: "simulated", Table: "events", OccurredAt: time.Now().UTC(), Position: transaction.Position{LSN: uint64(len(m.Payload))}}
			if err := emit(e); err != nil {
				return err
			}
		}
	}
}
