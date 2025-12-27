package ports

import (
	"context"
	"time"
)

// IKVCachePort interface for key-value cache.
type IKVCachePort interface {
	// Set stores a value in storage with the key and time-to-live (TTL).
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	// MSet is like Set but accepts multiple values:
	MSet(ctx context.Context, ttl time.Duration, kv map[string]interface{}) error
	// SetNX stores a value in storage with the key and time-to-live (TTL) if the key does not exist.
	SetNX(ctx context.Context, key string, value []byte, ttl time.Duration) error
	// Get retrieves a value from storage by key.
	Get(ctx context.Context, key string) ([]byte, error)
	// MGet
	MGet(ctx context.Context, key ...string) ([]any, error)
	// Exists checks if a key exists in storage.
	Exists(ctx context.Context, key string) (bool, error)
	// DelByPrefix deletes all keys that match the given prefix.
	DelByPrefix(ctx context.Context, prefix string) error
	// Del removes a specific key from storage.
	Del(ctx context.Context, key string) error
}

// INumericCache interface for numeric key-value cache.
type INumericCachePort interface {
	// Get retrieves a value from storage by key.
	Get(ctx context.Context, key string) (int64, error)
	// Set sets a value in the hash set.
	Set(ctx context.Context, key string, value int64, ttl time.Duration) error
	// SetNX sets a value in the hash set if the key does not exist.
	SetNX(ctx context.Context, key string, value int64, ttl time.Duration) error
	// Increment increments the value of a key by the given increment.
	Increment(ctx context.Context, key string, increment int64) (int64, error)
	// IncrementNX increments the value of a key by the given increment if the key does exist.
	IncrementEX(ctx context.Context, key string, increment int64) (int64, error)
	// Decrement decrements the value of a key by the given decrement.
	Decrement(ctx context.Context, key string, decrement int64) (int64, error)
	// DecrementNX decrements the value of a key by the given decrement if the key does exist.
	DecrementEX(ctx context.Context, key string, decrement int64) (int64, error)
	// Del removes a key from the hash set.
	Del(ctx context.Context, key string) error
	// ScanFunc
	ScanFunc(ctx context.Context, prefix string, scanFunc func(key string, v int64, err error)) error
}

// ISetPort interface for set port.
type ISetPort interface {
	// Add add value in to the Set.
	Add(ctx context.Context, key string, value string) error
	// IsMember checks if a value is a member of the Set.
	IsMember(ctx context.Context, key string, value string) (bool, error)
	// Members get all member in the Set.
	Members(ctx context.Context, key string) ([]string, error)
	// Del delete a value from the Set.
	Del(ctx context.Context, key string, value string) error
	// Pop get a random value from the Set and delete that value.
	Pop(ctx context.Context, key string) (string, error)
}

// ISetPort interface for set port.
// TODO: add more methods and change function arguments.
type IZSetPort interface {
	// Add add value in to the Set.
	Add(ctx context.Context, key string, value string) error
	// IsMember checks if a value is a member of the Set.
	IsMember(ctx context.Context, key string, value string) (bool, error)
	// Members get all member in the Set.
	Members(ctx context.Context, key string) ([]string, error)
	// Del delete a value from the Set.
	Del(ctx context.Context, key string, value string) error
	// Pop get a random value from the Set and delete that value.
	Pop(ctx context.Context, key string) (string, error)
}
