package bufferdomain

import (
	transaction "github.com/acme/streamforge-cdc/internal/transaction/domain"
	"sync"
)

type Record struct {
	ID       string                  `json:"id"`
	Tx       transaction.Transaction `json:"transaction"`
	Bytes    int                     `json:"bytes"`
	Attempts int                     `json:"attempts"`
	Dead     bool                    `json:"dead"`
}
type Queue struct {
	mu      sync.Mutex
	records []Record
	limit   int
	used    int
}

func New(limit int) *Queue {
	if limit < 1 {
		limit = 1024 * 1024
	}
	return &Queue{limit: limit}
}
func (q *Queue) Enqueue(tx transaction.Transaction) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	b := 0
	for _, e := range tx.Events {
		b += len(e.ID) + len(e.Table) + 128
	}
	id := tx.ID
	if len(tx.Events) > 0 {
		id = tx.Events[0].ID
	}
	q.records = append(q.records, Record{ID: id, Tx: tx, Bytes: b})
	q.used += b
	return nil
}
func (q *Queue) Peek() (Record, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.records) == 0 {
		return Record{}, false
	}
	return q.records[0], true
}
func (q *Queue) Ack(id string) bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.records) == 0 {
		return false
	}
	q.used -= q.records[0].Bytes
	q.records = q.records[1:]
	return true
}
func (q *Queue) Fail(id string, max int) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.records) == 0 || q.records[0].ID != id {
		return
	}
	q.records[0].Attempts++
	if q.records[0].Attempts >= max {
		q.records[0].Dead = true
	}
}
func (q *Queue) Stats() (int, int, int) {
	q.mu.Lock()
	defer q.mu.Unlock()
	dead := 0
	for _, r := range q.records {
		if r.Dead {
			dead++
		}
	}
	return len(q.records), q.used, dead
}
