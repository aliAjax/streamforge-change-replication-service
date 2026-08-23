package main

import (
	"context"
	"fmt"
	"github.com/acme/streamforge-cdc/internal/buffer/domain"
	"github.com/acme/streamforge-cdc/internal/delivery/application"
	"github.com/acme/streamforge-cdc/internal/delivery/domain"
	"github.com/acme/streamforge-cdc/internal/pipeline/domain"
	"github.com/acme/streamforge-cdc/internal/platform/config"
	"github.com/acme/streamforge-cdc/internal/platform/events"
	"github.com/acme/streamforge-cdc/internal/platform/httpapi"
	"github.com/acme/streamforge-cdc/internal/platform/logging"
	"github.com/acme/streamforge-cdc/internal/platform/metrics"
	"github.com/acme/streamforge-cdc/internal/source/application"
	sourcedomain "github.com/acme/streamforge-cdc/internal/source/domain"
	sourceinfra "github.com/acme/streamforge-cdc/internal/source/infrastructure"
	transaction "github.com/acme/streamforge-cdc/internal/transaction/domain"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type sink struct{ store *eventstore.Store }

func (s sink) Persist(ctx context.Context, t transaction.Transaction) error {
	return s.store.Persist(ctx, t)
}
func main() {
	cfg := config.Load()
	log := logging.New()
	repo := sourceinfra.NewMemoryRepository()
	store := eventstore.Store{}
	txsink := sink{store: &store}
	_ = txsink
	srcSvc := sourceapp.New(repo, sourceinfra.Validator{})
	pipes := pipelinedomain.NewRegistry()
	m := &metrics.Registry{}
	q := bufferdomain.New(cfg.BufferBytes)
	sender := &deliverydomain.MemorySender{}
	del := deliveryapp.Service{Queue: q, Sender: sender, MaxAttempts: 3, Batch: 10}
	_ = del
	_ = transaction.Operation("")
	_ = sourcedomain.Active
	server := &httpapi.Server{Sources: srcSvc, SourceRepo: repo, Pipelines: pipes, Events: &store, Metrics: m, APIKey: cfg.APIKey, Started: time.Now()}
	srv := &http.Server{Addr: cfg.HTTPAddr, Handler: server.Handler()}
	go func() {
		log.Info("streamforge started", map[string]any{"addr": cfg.HTTPAddr})
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("server stopped", map[string]any{"error": err})
		}
	}()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()
	shutdown, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(shutdown); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
