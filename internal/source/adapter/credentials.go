package sourceadapter

import (
	"context"
	"fmt"
	"sync"
)

type Credentials struct {
	Username string
	Password string
	TLS      bool
}
type Provider interface {
	Get(context.Context, string) (Credentials, error)
}
type MemoryProvider struct {
	mu    sync.RWMutex
	items map[string]Credentials
}

func NewProvider() *MemoryProvider { return &MemoryProvider{items: map[string]Credentials{}} }
func (p *MemoryProvider) Put(id string, c Credentials) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.items[id] = c
}
func (p *MemoryProvider) Get(_ context.Context, id string) (Credentials, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	c, ok := p.items[id]
	if !ok {
		return c, fmt.Errorf("credential %s not found", id)
	}
	return c, nil
}
