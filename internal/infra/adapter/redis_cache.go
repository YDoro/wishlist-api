package adapter

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/ydoro/wishlist/internal/domain"
)

type redisCache struct {
	client *redis.Client
}

type RedisCacheConfig struct {
	Addr     string
	Password string
	DB       int
}

func NewRedisCache(config RedisCacheConfig) domain.Cache {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
	})

	return &redisCache{
		client: client,
	}
}

func (c *redisCache) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	return c.client.Set(ctx, key, value, expiration).Err()
}

func (c *redisCache) Get(ctx context.Context, key string) (string, error) {
	result, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return result, err
}

func (c *redisCache) Delete(ctx context.Context, keys ...string) error {
	return c.client.Del(ctx, keys...).Err()
}
