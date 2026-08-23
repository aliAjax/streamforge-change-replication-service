package pipelinedomain

import (
	"fmt"
	source "github.com/acme/streamforge-cdc/internal/source/domain"
	"sync"
)

type Status string

const (
	PipelineDraft   Status = "draft"
	PipelineRunning Status = "running"
	PipelinePaused  Status = "paused"
	PipelineFailed  Status = "failed"
)

type Pipeline struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	SourceID    string `json:"source_id"`
	Strategy    string `json:"strategy"`
	Status      Status `json:"status"`
	Destination string `json:"destination"`
	Version     int64  `json:"version"`
}
type Registry struct {
	mu    sync.RWMutex
	items map[string]Pipeline
}

func NewRegistry() *Registry { return &Registry{items: map[string]Pipeline{}} }
func (r *Registry) Put(p Pipeline) error {
	if p.ID == "" || p.SourceID == "" {
		return fmt.Errorf("pipeline id and source required")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	p.Version++
	r.items[p.ID] = p
	return nil
}
func (r *Registry) Get(id string) (Pipeline, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.items[id]
	return p, ok
}
func (r *Registry) List() []Pipeline {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Pipeline, 0, len(r.items))
	for _, p := range r.items {
		out = append(out, p)
	}
	return out
}

var _ source.Kind
