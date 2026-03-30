package consumer

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/prometheus/client_golang/prometheus"
	"planet-scale-event-streaming/go-consumer/internal/db"
	"planet-scale-event-streaming/go-consumer/internal/model"
	"planet-scale-event-streaming/go-consumer/internal/redisx"
)

var (
	ProcessedCounter = prometheus.NewCounter(prometheus.CounterOpts{Name: "consumer_processed_total", Help: "Processed events"})
	FailedCounter    = prometheus.NewCounter(prometheus.CounterOpts{Name: "consumer_failed_total", Help: "Failed events"})
	LatencyHistogram = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "consumer_processing_seconds",
		Help:    "Event processing latency",
		Buckets: prometheus.DefBuckets,
	})
)

type Handler struct {
	store      *db.Store
	redis      *redisx.Client
	producer   sarama.SyncProducer
	retryTopic string
	dlqTopic   string
	maxRetries int
}

func NewHandler(store *db.Store, redis *redisx.Client, producer sarama.SyncProducer, retryTopic, dlqTopic string, maxRetries int) *Handler {
	return &Handler{store: store, redis: redis, producer: producer, retryTopic: retryTopic, dlqTopic: dlqTopic, maxRetries: maxRetries}
}

func (h *Handler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *Handler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *Handler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		started := time.Now()
		var event model.Event
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("unmarshal failed: %v", err)
			session.MarkMessage(msg, "bad-payload")
			continue
		}

		ctx := context.Background()
		seen, err := h.redis.Seen(ctx, event.EventID)
		if err != nil {
			log.Printf("redis check failed: %v", err)
		}
		if seen {
			log.Printf("duplicate skipped: %s", event.EventID)
			session.MarkMessage(msg, "duplicate")
			continue
		}

		if err := h.process(ctx, event); err != nil {
			FailedCounter.Inc()
			log.Printf("process failed for %s: %v", event.EventID, err)
			h.handleFailure(event)
			session.MarkMessage(msg, "failed-forwarded")
			continue
		}

		if err := h.store.SaveProcessed(ctx, event, "processed"); err != nil {
			log.Printf("db save failed: %v", err)
		}
		ProcessedCounter.Inc()
		LatencyHistogram.Observe(time.Since(started).Seconds())
		session.MarkMessage(msg, "processed")
	}
	return nil
}

func (h *Handler) process(ctx context.Context, event model.Event) error {
	_ = ctx
	if simulate, ok := event.Payload["simulateFailure"].(bool); ok && simulate {
		return sarama.ErrOutOfBrokers
	}
	log.Printf("processed event=%s type=%s", event.EventID, event.EventType)
	return nil
}

func (h *Handler) handleFailure(event model.Event) {
	event.RetryCount++
	topic := h.retryTopic
	if event.RetryCount > h.maxRetries {
		topic = h.dlqTopic
	}
	payload, _ := json.Marshal(event)
	_, _, err := h.producer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(event.PartitionKey),
		Value: sarama.ByteEncoder(payload),
	})
	if err != nil {
		log.Printf("failed to publish retry/dlq event: %v", err)
	}
}
