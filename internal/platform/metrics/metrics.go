package metrics

import (
	"fmt"
	"sync/atomic"
)

type Registry struct {
	requests atomic.Uint64
	events   atomic.Uint64
	failures atomic.Uint64
}

func (r *Registry) Request() { r.requests.Add(1) }
func (r *Registry) Event()   { r.events.Add(1) }
func (r *Registry) Failure() { r.failures.Add(1) }
func (r *Registry) Prometheus() string {
	return fmt.Sprintf("# TYPE streamforge_requests_total counter\nstreamforge_requests_total %d\n# TYPE streamforge_events_total counter\nstreamforge_events_total %d\n# TYPE streamforge_failures_total counter\nstreamforge_failures_total %d\n", r.requests.Load(), r.events.Load(), r.failures.Load())
}
