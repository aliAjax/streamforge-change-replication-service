package mysqladapter

import (
	"encoding/binary"
	"fmt"
	transaction "github.com/acme/streamforge-cdc/internal/transaction/domain"
)

const (
	EventFormatDescription byte = 0x0f
	EventTableMap          byte = 0x13
	EventWriteRows         byte = 0x1e
	EventUpdateRows        byte = 0x1f
	EventDeleteRows        byte = 0x20
	EventRotate            byte = 0x04
)

func EncodeEvent(kind byte, pos uint32, payload []byte) []byte {
	h := make([]byte, 13)
	binary.LittleEndian.PutUint32(h, uint32(1))
	h[4] = kind
	binary.LittleEndian.PutUint32(h[5:], pos)
	binary.LittleEndian.PutUint32(h[9:], uint32(len(payload)))
	return append(h, payload...)
}
func DecodeRow(payload []byte, columns int) ([]any, error) {
	if columns < 0 || columns > 1024 {
		return nil, fmt.Errorf("invalid column count")
	}
	row := make([]any, columns)
	for i := range row {
		if i < len(payload) {
			row[i] = payload[i]
		} else {
			row[i] = nil
		}
	}
	return row, nil
}
func Position(kind byte, pos uint32) transaction.Position {
	return transaction.Position{File: fmt.Sprintf("binlog.%06d", kind), Offset: uint64(pos)}
}
