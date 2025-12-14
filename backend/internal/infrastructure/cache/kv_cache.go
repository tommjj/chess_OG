package cache

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/tommjj/chess_OG/backend/internal/core/domain"
)

type kvcache struct {
	redis *Redis
}

func NewKeyValueCacheAdapter(redis *Redis) *kvcache {
	return &kvcache{
		redis: redis,
	}
}

func (r *kvcache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return r.redis.Set(ctx, key, value, ttl).Err()
}

// MSet is like Set but accepts multiple values:
//
//	MSet("key1", "value1", "key2", "value2")
func (r *kvcache) MSet(ctx context.Context, ttl time.Duration, kv map[string]interface{}) error {
	pipe := r.redis.Pipeline()

	pipe.MSet(ctx, kv)

	if ttl > 0 {
		for key := range kv {
			pipe.Expire(ctx, key, ttl)
		}
	}

	_, err := pipe.Exec(ctx)
	return err
}

func (r *kvcache) SetNX(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return r.redis.SetEx(ctx, key, value, ttl).Err()
}

func (r *kvcache) Get(ctx context.Context, key string) ([]byte, error) {
	val, err := r.redis.Get(ctx, key).Result()
	if err != nil {
		return nil, handleRedisErr(err)
	}

	return []byte(val), nil
}

func (r *kvcache) MGet(ctx context.Context, key ...string) ([]any, error) {
	vals, err := r.redis.MGet(ctx, key...).Result()
	if err != nil {
		return nil, handleRedisErr(err)
	}

	return vals, nil
}

func (r *kvcache) Exists(ctx context.Context, key string) (bool, error) {
	exists, err := r.redis.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return exists > 0, nil
}

func (r *kvcache) DelByPrefix(ctx context.Context, prefix string) error {
	var cursor uint64
	var keys []string

	for {
		var err error
		keys, cursor, err = r.redis.Scan(ctx, cursor, prefix+"*", 100).Result()
		if err != nil {
			return err
		}

		for _, key := range keys {
			err := r.redis.Del(ctx, key).Err()
			if err != nil {
				return err
			}
		}

		if cursor == 0 {
			break
		}
	}

	return nil
}

func (r *kvcache) Del(ctx context.Context, key string) error {
	return r.redis.Del(ctx, key).Err()
}

func handleRedisErr(err error) error {
	if errors.Is(err, redis.Nil) {
		return domain.ErrDataNotFound
	}
	return err
}
