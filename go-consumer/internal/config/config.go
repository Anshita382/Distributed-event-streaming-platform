package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	KafkaBrokers []string
	GroupID      string
	MainTopic    string
	RetryTopic   string
	DLQTopic     string
	RedisAddr    string
	PostgresURL  string
	MetricsAddr  string
	MaxRetries   int
}

func Load() Config {
	return Config{
		KafkaBrokers: splitCSV(getEnv("GO_CONSUMER_KAFKA_BROKERS", "localhost:9092")),
		GroupID:      getEnv("GO_CONSUMER_GROUP_ID", "go-consumer-group"),
		MainTopic:    getEnv("GO_CONSUMER_MAIN_TOPIC", "events.main"),
		RetryTopic:   getEnv("GO_CONSUMER_RETRY_TOPIC", "events.retry"),
		DLQTopic:     getEnv("GO_CONSUMER_DLQ_TOPIC", "events.dlq"),
		RedisAddr:    getEnv("GO_CONSUMER_REDIS_ADDR", "localhost:6379"),
		PostgresURL:  getEnv("GO_CONSUMER_POSTGRES_URL", "postgres://postgres:postgres@localhost:5432/eventing?sslmode=disable"),
		MetricsAddr:  getEnv("GO_CONSUMER_METRICS_ADDR", ":2112"),
		MaxRetries:   getEnvInt("GO_CONSUMER_MAX_RETRIES", 3),
	}
}

func getEnv(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		parsed, err := strconv.Atoi(value)
		if err == nil {
			return parsed
		}
	}
	return fallback
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	if len(result) == 0 {
		return []string{"localhost:9092"}
	}
	return result
}
