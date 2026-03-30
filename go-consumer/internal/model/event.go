package model

import "time"

type Event struct {
	EventID      string         `json:"eventId"`
	EventType    string         `json:"eventType"`
	PartitionKey string         `json:"partitionKey"`
	Payload      map[string]any `json:"payload"`
	RetryCount   int            `json:"retryCount"`
	CreatedAt    time.Time      `json:"createdAt"`
}
