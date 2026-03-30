package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	KafkaBrokers []string
	Topic        string
	MetricsAddr  string
	BurstSize    int
}

func Load() Config {
	return Config{
		KafkaBrokers: splitCSV(getEnv("GO_PRODUCER_KAFKA_BROKERS", "localhost:9092")),
		Topic:        getEnv("GO_PRODUCER_TOPIC", "events.main"),
		MetricsAddr:  getEnv("GO_PRODUCER_METRICS_ADDR", ":2113"),
		BurstSize:    getEnvInt("GO_PRODUCER_BURST_SIZE", 500),
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
