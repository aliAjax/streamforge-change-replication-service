package deliveryadapter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	delivery "github.com/acme/streamforge-cdc/internal/delivery/domain"
	"net/http"
	"time"
)

type HTTPSender struct {
	URL    string
	Client *http.Client
}

func (h HTTPSender) Send(ctx context.Context, msgs []delivery.Message) error {
	if h.URL == "" {
		return fmt.Errorf("webhook URL is empty")
	}
	if h.Client == nil {
		h.Client = &http.Client{Timeout: 5 * time.Second}
	}
	b, err := json.Marshal(msgs)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, h.URL, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := h.Client.Do(req)
	if err != nil {
		return fmt.Errorf("webhook: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("webhook status %s", resp.Status)
	}
	return nil
}
