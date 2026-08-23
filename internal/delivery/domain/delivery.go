package deliverydomain

import (
	"context"
	"fmt"
	transaction "github.com/acme/streamforge-cdc/internal/transaction/domain"
	"sync"
)

type Message struct {
	ID          string                  `json:"id"`
	Transaction transaction.Transaction `json:"transaction"`
}
type Sender interface {
	Send(context.Context, []Message) error
}
type MemorySender struct {
	mu       sync.Mutex
	Messages []Message
	Fail     bool
}

func (m *MemorySender) Send(_ context.Context, msgs []Message) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.Fail {
		return fmt.Errorf("simulated downstream failure")
	}
	m.Messages = append(m.Messages, msgs...)
	return nil
}
func (m *MemorySender) Count() int { m.mu.Lock(); defer m.mu.Unlock(); return len(m.Messages) }

type WebhookSender struct {
	URL    string
	Client interface{}
}
