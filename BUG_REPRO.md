# Bug Reproduction

## Bug

When a handler writes a successful status/body and then panics, the recovery middleware writes a second 500 response to the same `http.ResponseWriter`. The committed response status remains 200 while the error JSON is appended to the already-written body.

## Trigger

Run a handler that writes a partial response and panics before returning through `internal/platform/httpx.Middleware.Wrap`.

## Observed error

The real reproduction logged `http: superfluous response.WriteHeader call from github.com/acme/streamforge-cdc/internal/platform/httpx.JSON (httpx.go:39)` and returned `STATUS: 200` with a body containing both the partial success JSON and `{"error":"internal server error","status":500}`.
