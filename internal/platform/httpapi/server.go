package httpapi

import (
	"encoding/json"
	"fmt"
	pipeline "github.com/acme/streamforge-cdc/internal/pipeline/domain"
	eventstore "github.com/acme/streamforge-cdc/internal/platform/events"
	"github.com/acme/streamforge-cdc/internal/platform/httpx"
	"github.com/acme/streamforge-cdc/internal/platform/metrics"
	sourceapp "github.com/acme/streamforge-cdc/internal/source/application"
	sourcedomain "github.com/acme/streamforge-cdc/internal/source/domain"
	transaction "github.com/acme/streamforge-cdc/internal/transaction/domain"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Server struct {
	Sources    *sourceapp.Service
	SourceRepo sourceapp.Repository
	Pipelines  *pipeline.Registry
	Events     *eventstore.Store
	Metrics    *metrics.Registry
	APIKey     string
	Started    time.Time
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.health)
	mux.HandleFunc("/readyz", s.ready)
	mux.HandleFunc("/metrics", s.metrics)
	mux.HandleFunc("/api/v1/sources", s.sources)
	mux.HandleFunc("/api/v1/pipelines", s.pipelines)
	mux.HandleFunc("/api/v1/events", s.events)
	mux.HandleFunc("/api/v1/simulate/transaction", s.simulate)
	return httpx.Middleware{APIKey: s.APIKey, Requests: s.Metrics.Request}.Wrap(mux)
}
func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	httpx.JSON(w, 200, map[string]any{"status": "ok", "uptime": time.Since(s.Started).String()})
}
func (s *Server) ready(w http.ResponseWriter, _ *http.Request) {
	if s.SourceRepo == nil {
		httpx.Error(w, 503, fmt.Errorf("source repository unavailable"))
		return
	}
	httpx.JSON(w, 200, map[string]string{"status": "ready"})
}
func (s *Server) metrics(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	_, _ = w.Write([]byte(s.Metrics.Prometheus()))
}
func (s *Server) sources(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		items, err := s.Sources.List(r.Context())
		if err != nil {
			httpx.Error(w, 500, err)
			return
		}
		httpx.JSON(w, 200, map[string]any{"items": items})
	case http.MethodPost:
		var in sourcedomain.Source
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			httpx.Error(w, 400, err)
			return
		}
		out, err := s.Sources.Create(r.Context(), in)
		if err != nil {
			httpx.Error(w, 400, err)
			return
		}
		httpx.JSON(w, 201, out)
	default:
		http.NotFound(w, r)
	}
}
func (s *Server) pipelines(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		httpx.JSON(w, 200, map[string]any{"items": s.Pipelines.List()})
		return
	}
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}
	var in pipeline.Pipeline
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		httpx.Error(w, 400, err)
		return
	}
	if err := s.Pipelines.Put(in); err != nil {
		httpx.Error(w, 400, err)
		return
	}
	httpx.JSON(w, 201, in)
}
func (s *Server) events(w http.ResponseWriter, r *http.Request) {
	cursor, _ := strconv.Atoi(r.URL.Query().Get("cursor"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit == 0 {
		limit = 100
	}
	items := s.Events.List(cursor, limit)
	httpx.JSON(w, 200, map[string]any{"items": items, "next_cursor": cursor + len(items), "count": len(items)})
}

func (s *Server) simulate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}
	var in struct {
		Source, Schema, Table, TransactionID string
		Operation                            transaction.Operation
		Before, After, PrimaryKey            map[string]any
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		httpx.Error(w, http.StatusBadRequest, err)
		return
	}
	if in.Source == "" {
		in.Source = "sim-postgres"
	}
	if in.Schema == "" {
		in.Schema = "public"
	}
	if in.Table == "" {
		in.Table = "events"
	}
	if in.TransactionID == "" {
		in.TransactionID = fmt.Sprintf("tx-%d", time.Now().UnixNano())
	}
	if in.Operation == "" {
		in.Operation = transaction.Insert
	}
	e := transaction.Event{Source: in.Source, Schema: in.Schema, Table: in.Table, Operation: in.Operation, Before: in.Before, After: in.After, PrimaryKey: in.PrimaryKey, TransactionID: in.TransactionID, Position: transaction.Position{LSN: uint64(s.Events.Count() + 1)}, OccurredAt: time.Now().UTC(), SchemaVersion: 1, Sequence: 1}
	e.ID = e.ComputeID()
	t := transaction.Transaction{ID: in.TransactionID, Source: in.Source, Position: e.Position, StartedAt: e.OccurredAt, CommittedAt: e.OccurredAt, Events: []transaction.Event{e}}
	if err := s.Events.Persist(r.Context(), t); err != nil {
		httpx.Error(w, http.StatusConflict, err)
		return
	}
	s.Metrics.Event()
	httpx.JSON(w, http.StatusAccepted, map[string]any{"event": e, "checkpoint": e.Position})
}
func ParseBearer(v string) string { return strings.TrimPrefix(v, "Bearer ") }
