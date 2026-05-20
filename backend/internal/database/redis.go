package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sulthancom/pos-backend/internal/config"
)

type Redis struct {
	Client *redis.Client
}

func NewRedis(cfg *config.Config) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Printf("WARNING: Redis not available, continuing without cache: %v", err)
		return &Redis{Client: client}, nil
	}

	log.Printf("Connected to Redis: %s", cfg.RedisAddr)
	return &Redis{Client: client}, nil
}

func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}

func (r *Redis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.Client.Set(ctx, key, value, expiration).Err()
}

func (r *Redis) Delete(ctx context.Context, keys ...string) error {
	return r.Client.Del(ctx, keys...).Err()
}

func (r *Redis) Exists(ctx context.Context, key string) (bool, error) {
	result, err := r.Client.Exists(ctx, key).Result()
	return result > 0, err
}

func (r *Redis) Close() error {
	return r.Client.Close()
}

func GenerateCacheKey(prefix string, params map[string]string) string {
	key := prefix
	for k, v := range params {
		key += fmt.Sprintf(":%s=%s", k, v)
	}
	return key
}
