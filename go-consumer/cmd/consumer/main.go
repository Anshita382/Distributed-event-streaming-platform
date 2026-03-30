package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"planet-scale-event-streaming/go-consumer/internal/config"
	"planet-scale-event-streaming/go-consumer/internal/consumer"
	"planet-scale-event-streaming/go-consumer/internal/db"
	"planet-scale-event-streaming/go-consumer/internal/redisx"
)

func main() {
	cfg := config.Load()
	prometheus.MustRegister(consumer.ProcessedCounter, consumer.FailedCounter, consumer.LatencyHistogram)

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Printf("consumer metrics listening on %s", cfg.MetricsAddr)
		if err := http.ListenAndServe(cfg.MetricsAddr, nil); err != nil {
			log.Fatal(err)
		}
	}()

	ctx := context.Background()
	store, err := db.New(ctx, cfg.PostgresURL)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Pool.Close()

	redisClient := redisx.New(cfg.RedisAddr)

	saramaCfg := sarama.NewConfig()
	saramaCfg.Version = sarama.V2_8_0_0
	saramaCfg.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRange()
	saramaCfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	saramaCfg.Consumer.Fetch.Default = 1024 * 1024
	saramaCfg.Consumer.MaxProcessingTime = 500000000
	saramaCfg.ChannelBufferSize = 256
	producerCfg := sarama.NewConfig()
	producerCfg.Version = sarama.V2_8_0_0
	producerCfg.Producer.Return.Successes = true
	producerCfg.Producer.RequiredAcks = sarama.WaitForAll
	producerCfg.Producer.Retry.Max = 5
	producerCfg.Producer.Idempotent = true
	producerCfg.Net.MaxOpenRequests = 1

	retryProducer, err := sarama.NewSyncProducer(cfg.KafkaBrokers, producerCfg)
	if err != nil {
		log.Fatal(err)
	}
	defer retryProducer.Close()

	handler := consumer.NewHandler(store, redisClient, retryProducer, cfg.RetryTopic, cfg.DLQTopic, cfg.MaxRetries)

	group, err := sarama.NewConsumerGroup(cfg.KafkaBrokers, cfg.GroupID, saramaCfg)
	if err != nil {
		log.Fatal(err)
	}
	defer group.Close()

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			if err := group.Consume(ctx, []string{cfg.MainTopic, cfg.RetryTopic}, handler); err != nil {
				log.Printf("consume error: %v", err)
			}
			if ctx.Err() != nil {
				return
			}
		}
	}()

	log.Println("go consumer started")
	<-sigterm
	log.Println("shutting down consumer")
}
