package transactiondomain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

type Operation string

const (
	Insert   Operation = "insert"
	Update   Operation = "update"
	Delete   Operation = "delete"
	Truncate Operation = "truncate"
	DDL      Operation = "ddl"
)

type Position struct {
	LSN    uint64 `json:"lsn,omitempty"`
	GTID   string `json:"gtid,omitempty"`
	File   string `json:"file,omitempty"`
	Offset uint64 `json:"offset,omitempty"`
}

func (p Position) String() string {
	return fmt.Sprintf("lsn=%d gtid=%s file=%s:%d", p.LSN, p.GTID, p.File, p.Offset)
}

type Event struct {
	ID            string            `json:"id"`
	Source        string            `json:"source"`
	Database      string            `json:"database"`
	Schema        string            `json:"schema"`
	Table         string            `json:"table"`
	Operation     Operation         `json:"operation"`
	Before        map[string]any    `json:"before,omitempty"`
	After         map[string]any    `json:"after,omitempty"`
	PrimaryKey    map[string]any    `json:"primary_key,omitempty"`
	TransactionID string            `json:"transaction_id"`
	Position      Position          `json:"position"`
	OccurredAt    time.Time         `json:"occurred_at"`
	SchemaVersion int64             `json:"schema_version"`
	Sequence      int               `json:"sequence"`
	Metadata      map[string]string `json:"metadata,omitempty"`
}

func (e Event) ComputeID() string {
	b, _ := json.Marshal(struct {
		Source string
		Tx     string
		Pos    Position
		Seq    int
	}{e.Source, e.TransactionID, e.Position, e.Sequence})
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:])
}

type Transaction struct {
	ID          string    `json:"id"`
	Source      string    `json:"source"`
	Position    Position  `json:"position"`
	StartedAt   time.Time `json:"started_at"`
	CommittedAt time.Time `json:"committed_at"`
	Events      []Event   `json:"events"`
}

func (t *Transaction) Append(e Event) error {
	if e.TransactionID != t.ID {
		return fmt.Errorf("transaction %s: event belongs to %s", t.ID, e.TransactionID)
	}
	e.Sequence = len(t.Events) + 1
	if e.ID == "" {
		e.ID = e.ComputeID()
	}
	t.Events = append(t.Events, e)
	return nil
}

func (t *Transaction) Commit(at time.Time) error {
	if len(t.Events) == 0 {
		return fmt.Errorf("transaction %s: empty transaction", t.ID)
	}
	if at.Before(t.StartedAt) {
		return fmt.Errorf("transaction %s: commit precedes start", t.ID)
	}
	t.CommittedAt = at
	return nil
}
