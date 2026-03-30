package db

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"planet-scale-event-streaming/go-consumer/internal/model"
)

type Store struct {
	Pool *pgxpool.Pool
}

func New(ctx context.Context, postgresURL string) (*Store, error) {
	pool, err := pgxpool.New(ctx, postgresURL)
	if err != nil {
		return nil, err
	}
	return &Store{Pool: pool}, nil
}

func (s *Store) SaveProcessed(ctx context.Context, event model.Event, status string) error {
	payload, err := json.Marshal(event.Payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}
	_, err = s.Pool.Exec(ctx, `
		INSERT INTO processed_events(event_id, event_type, partition_key, payload, status, retry_count)
		VALUES ($1, $2, $3, $4::jsonb, $5, $6)
		ON CONFLICT (event_id) DO NOTHING
	`, event.EventID, event.EventType, event.PartitionKey, string(payload), status, event.RetryCount)
	return err
}
