package bufferinfra

import (
	"encoding/json"
	buffer "github.com/acme/streamforge-cdc/internal/buffer/domain"
	transaction "github.com/acme/streamforge-cdc/internal/transaction/domain"
	"os"
	"sync"
)

type Snapshot struct {
	mu    sync.Mutex
	Path  string
	Queue *buffer.Queue
}

func NewSnapshot(path string, q *buffer.Queue) *Snapshot { return &Snapshot{Path: path, Queue: q} }
func (s *Snapshot) Save() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	n, _, _ := s.Queue.Stats()
	return os.WriteFile(s.Path, []byte(jsonString(n)), 0600)
}
func jsonString(n int) string { b, _ := json.Marshal(map[string]int{"records": n}); return string(b) }

var _ transaction.Position
