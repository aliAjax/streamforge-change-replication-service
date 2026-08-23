package deliveryadapter

import (
	"context"
	"fmt"
	delivery "github.com/acme/streamforge-cdc/internal/delivery/domain"
	"sync"
)

type KafkaMessage struct {
	Topic string
	Key   string
	Value []byte
}
type KafkaProducer interface {
	Produce(context.Context, KafkaMessage) error
}
type MemoryKafka struct {
	mu       sync.Mutex
	Messages []KafkaMessage
	Fail     bool
}

func (k *MemoryKafka) Produce(_ context.Context, m KafkaMessage) error {
	k.mu.Lock()
	defer k.mu.Unlock()
	if k.Fail {
		return fmt.Errorf("kafka unavailable")
	}
	k.Messages = append(k.Messages, m)
	return nil
}

type KafkaSender struct {
	Producer KafkaProducer
	Topic    string
}

func (s KafkaSender) Send(ctx context.Context, msgs []delivery.Message) error {
	for _, m := range msgs {
		if err := s.Producer.Produce(ctx, KafkaMessage{Topic: s.Topic, Key: m.ID, Value: []byte(m.ID)}); err != nil {
			return fmt.Errorf("kafka produce: %w", err)
		}
	}
	return nil
}
