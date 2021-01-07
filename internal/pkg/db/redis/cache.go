package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// Cache struct
type Cache struct {
	client redis.UniversalClient
	ttl    time.Duration
}

const apqPrefix = ""

// NewCache to create new object Cache
func NewCache(ctx context.Context, redisAddress string, password string, ttl time.Duration) (*Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Password: password, // no password set
		DB:       0,        // use default DB
	})

	err := client.Ping(ctx).Err()
	if err != nil {
		return nil, fmt.Errorf("could not create cache: %w", err)
	}

	return &Cache{client: client, ttl: ttl}, nil
}

// SetTTL cache
func (c *Cache) SetTTL(ttl time.Duration) {
	c.ttl = ttl
}

// ResetTTL cache
func (c *Cache) ResetTTL() {
	c.ttl = 24 * time.Hour
}

// Add cache
func (c *Cache) Add(ctx context.Context, key string, value interface{}) {
	c.client.Set(ctx, apqPrefix+key, value, c.ttl)
}

// Get Cache
func (c *Cache) Get(ctx context.Context, key string) (interface{}, bool) {
	s, err := c.client.Get(ctx, apqPrefix+key).Result()
	if err != nil {
		return struct{}{}, false
	}
	return s, true
}

// DeleteByPrefix cache
func (c *Cache) DeleteByPrefix(ctx context.Context, prefix string) error {
	var err error
	iter := c.client.Scan(ctx, 0, prefix+"*", 0).Iterator()
	for iter.Next(ctx) {
		err = c.Del(ctx, iter.Val())
		if err != nil {
			return err
		}
	}

	return nil
}

// Del cache
func (c *Cache) Del(ctx context.Context, keys ...string) error {
	return c.client.Del(ctx, keys...).Err()
}
