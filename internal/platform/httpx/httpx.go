package httpx

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Middleware struct {
	APIKey   string
	Requests func()
}

func (m Middleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if m.Requests != nil {
			m.Requests()
		}
		ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
		defer cancel()
		r = r.WithContext(ctx)
		w.Header().Set("X-Request-ID", fmt.Sprintf("req-%d", time.Now().UnixNano()))
		if m.APIKey != "" && r.URL.Path != "/healthz" && r.URL.Path != "/readyz" && r.URL.Path != "/metrics" && r.Header.Get("Authorization") != "Bearer "+m.APIKey {
			writeError(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		defer func() {
			if recover() != nil {
				writeError(w, 500, "internal server error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}
func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func writeError(w http.ResponseWriter, status int, msg string) {
	JSON(w, status, map[string]any{"error": msg, "status": status})
}
func Error(w http.ResponseWriter, status int, err error) { writeError(w, status, err.Error()) }
