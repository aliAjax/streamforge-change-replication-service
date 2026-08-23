package pipelineadapter

import (
	"crypto/sha256"
	"fmt"
	transaction "github.com/acme/streamforge-cdc/internal/transaction/domain"
)

type Partitioner struct{ Partitions int }

func (p Partitioner) Partition(e transaction.Event) int {
	if p.Partitions < 1 {
		return 0
	}
	key := e.ID
	if len(e.PrimaryKey) > 0 {
		for k, v := range e.PrimaryKey {
			key += k + toString(v)
		}
	}
	h := sha256.Sum256([]byte(key))
	n := 0
	for _, b := range h[:] {
		n = (n*31 + int(b)) % p.Partitions
	}
	return n
}
func toString(v any) string { return fmt.Sprint(v) }
