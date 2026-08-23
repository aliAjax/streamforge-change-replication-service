package mysqldomain

import (
	"encoding/binary"
	"fmt"
	transaction "github.com/acme/streamforge-cdc/internal/transaction/domain"
	"io"
)

type Event struct {
	Type     byte
	ServerID uint32
	LogPos   uint32
	Payload  []byte
}

func Decode(r io.Reader, max int) (Event, error) {
	h := make([]byte, 13)
	if _, err := io.ReadFull(r, h); err != nil {
		return Event{}, err
	}
	n := int(binary.LittleEndian.Uint32(h[9:]))
	if max > 0 && n > max {
		return Event{}, fmt.Errorf("binlog event too large: %d", n)
	}
	p := make([]byte, n)
	if _, err := io.ReadFull(r, p); err != nil {
		return Event{}, err
	}
	return Event{Type: h[4], ServerID: binary.LittleEndian.Uint32(h), LogPos: binary.LittleEndian.Uint32(h[5:9]), Payload: p}, nil
}
func Position(e Event) transaction.Position {
	return transaction.Position{File: "binlog.000001", Offset: uint64(e.LogPos)}
}
