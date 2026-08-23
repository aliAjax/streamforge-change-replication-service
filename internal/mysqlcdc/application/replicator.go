package mysqlapp

import (
	"context"
	"fmt"
	mysql "github.com/acme/streamforge-cdc/internal/mysqlcdc/domain"
	transaction "github.com/acme/streamforge-cdc/internal/transaction/domain"
	"io"
	"time"
)

type Stream interface {
	Read(context.Context) (io.Reader, error)
	Confirm(context.Context, transaction.Position) error
}

func Run(ctx context.Context, s Stream, emit func(transaction.Event) error) error {
	r, err := s.Read(ctx)
	if err != nil {
		return fmt.Errorf("mysql stream: %w", err)
	}
	for {
		e, err := mysql.Decode(r, 16<<20)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("decode binlog: %w", err)
		}
		ev := transaction.Event{Source: "mysql", Operation: transaction.Update, TransactionID: "simulated", Schema: "public", Table: "events", OccurredAt: time.Now().UTC(), Position: mysql.Position(e)}
		if err := emit(ev); err != nil {
			return err
		}
	}
}
