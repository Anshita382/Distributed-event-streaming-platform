package redisx

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	inner *redis.Client
}

func New(addr string) *Client {
	return &Client{inner: redis.NewClient(&redis.Options{Addr: addr})}
}

func (c *Client) Seen(ctx context.Context, eventID string) (bool, error) {
	res, err := c.inner.SetNX(ctx, "evt:"+eventID, "1", 24*time.Hour).Result()
	if err != nil {
		return false, err
	}
	return !res, nil
}
