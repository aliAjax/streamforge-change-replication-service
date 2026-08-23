package postgresdomain

import (
	"encoding/binary"
	"fmt"
	transaction "github.com/acme/streamforge-cdc/internal/transaction/domain"
	"io"
)

type Message struct {
	Tag     byte
	Payload []byte
}

func Decode(r io.Reader, max int) (Message, error) {
	h := make([]byte, 5)
	if _, err := io.ReadFull(r, h[:1]); err != nil {
		return Message{}, err
	}
	if _, err := io.ReadFull(r, h[1:]); err != nil {
		return Message{}, err
	}
	n := int(binary.BigEndian.Uint32(h[1:]))
	if n < 4 || n > max {
		return Message{}, fmt.Errorf("postgres message length %d out of bounds", n)
	}
	p := make([]byte, n-4)
	if _, err := io.ReadFull(r, p); err != nil {
		return Message{}, fmt.Errorf("postgres payload: %w", err)
	}
	return Message{Tag: h[0], Payload: p}, nil
}

type Position struct{ LSN uint64 }

func PositionFromLSN(v uint64) transaction.Position { return transaction.Position{LSN: v} }
