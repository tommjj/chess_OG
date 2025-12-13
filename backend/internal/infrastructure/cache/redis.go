package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	*redis.Client
}

// NewRedis create new redis client
func NewRedis(addr, password string) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return &Redis{client}, nil
}

// Close close redis client
func (r *Redis) Close() error {
	return r.Client.Close()
}
