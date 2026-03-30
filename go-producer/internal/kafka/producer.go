package kafka

import (
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
	"planet-scale-event-streaming/go-producer/internal/model"
)

type Producer struct {
	syncProducer sarama.SyncProducer
	topic        string
}

func NewProducer(brokers []string, topic string) (*Producer, error) {
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V2_8_0_0
	cfg.Producer.Return.Successes = true
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Retry.Max = 10
	cfg.Producer.Idempotent = true
	cfg.Net.MaxOpenRequests = 1
	cfg.Producer.Compression = sarama.CompressionSnappy

	sp, err := sarama.NewSyncProducer(brokers, cfg)
	if err != nil {
		return nil, fmt.Errorf("create producer: %w", err)
	}
	return &Producer{syncProducer: sp, topic: topic}, nil
}

func (p *Producer) Publish(event model.Event) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}
	_, _, err = p.syncProducer.SendMessage(&sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.StringEncoder(event.PartitionKey),
		Value: sarama.ByteEncoder(payload),
	})
	if err != nil {
		return fmt.Errorf("publish event: %w", err)
	}
	return nil
}

func (p *Producer) Close() error {
	return p.syncProducer.Close()
}
