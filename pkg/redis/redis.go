package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const timeout = 5 * time.Second

func NewClient(redisURL string) (client *redis.Client, err error) {
	options, err := redis.ParseURL(redisURL)
	if err != nil {
		err = fmt.Errorf("parse url: %v", err)
		return
	}
	client = redis.NewClient(options)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	_, err = client.Ping(ctx).Result()
	if err != nil {
		err = fmt.Errorf("ping: %w", err)
		return
	}
	return
}
