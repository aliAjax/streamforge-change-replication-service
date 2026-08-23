package qa11

import (
	httpx "github.com/acme/streamforge-cdc/internal/platform/httpx"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMiddlewareRecoversWithoutLeakingPartialResponse(t *testing.T) {
	h := httpx.Middleware{}.Wrap(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(201)
		_, _ = w.Write([]byte("partial"))
		panic("boom")
	}))
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/x", nil))
	if rr.Code != 500 || !strings.Contains(rr.Body.String(), "internal server error") || strings.Contains(rr.Body.String(), "partial") {
		t.Fatalf("status=%d body=%q", rr.Code, rr.Body.String())
	}
}
