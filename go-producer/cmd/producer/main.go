package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"planet-scale-event-streaming/go-producer/internal/config"
	kafkax "planet-scale-event-streaming/go-producer/internal/kafka"
	"planet-scale-event-streaming/go-producer/internal/model"
)

var publishedCounter = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "go_producer_published_total",
	Help: "Total events published by Go producer",
})

func main() {
	cfg := config.Load()
	prometheus.MustRegister(publishedCounter)

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Printf("producer metrics listening on %s", cfg.MetricsAddr)
		if err := http.ListenAndServe(cfg.MetricsAddr, nil); err != nil {
			log.Fatal(err)
		}
	}()

	producer, err := kafkax.NewProducer(cfg.KafkaBrokers, cfg.Topic)
	if err != nil {
		log.Fatal(err)
	}
	defer producer.Close()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	log.Println("go producer started")
	for range ticker.C {
		for i := 0; i < cfg.BurstSize; i++ {
			evt := model.Event{
				EventID:      uuid.NewString(),
				EventType:    []string{"order.created", "inventory.updated", "payment.captured"}[rand.Intn(3)],
				PartitionKey: "tenant-" + string(rune('A'+rune(rand.Intn(5)))),
				Payload: map[string]any{
					"source":          "go-producer",
					"amount":          rand.Intn(10000),
					"simulateFailure": rand.Intn(500) == 1,
				},
				RetryCount: 0,
				CreatedAt:  time.Now().UTC(),
			}
			if err := producer.Publish(evt); err != nil {
				log.Printf("publish failed: %v", err)
				continue
			}
			publishedCounter.Inc()
		}
		log.Printf("published %d events this tick", cfg.BurstSize)
	}
}
