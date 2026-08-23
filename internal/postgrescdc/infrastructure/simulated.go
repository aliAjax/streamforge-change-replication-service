package postgresinfra

import (
	"bytes"
	"context"
	"encoding/binary"
	transaction "github.com/acme/streamforge-cdc/internal/transaction/domain"
	"io"
	"sync"
)

type SimulatedStream struct {
	mu        sync.Mutex
	Data      []byte
	Confirmed []transaction.Position
}

func NewSimulated() *SimulatedStream { return &SimulatedStream{} }
func (s *SimulatedStream) Push(tag byte, payload []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	h := make([]byte, 5)
	h[0] = tag
	binary.BigEndian.PutUint32(h[1:], uint32(len(payload)+4))
	s.Data = append(s.Data, h...)
	s.Data = append(s.Data, payload...)
}
func (s *SimulatedStream) Read(context.Context) (io.Reader, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return bytes.NewReader(append([]byte(nil), s.Data...)), nil
}
func (s *SimulatedStream) Confirm(_ context.Context, p transaction.Position) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Confirmed = append(s.Confirmed, p)
	return nil
}
