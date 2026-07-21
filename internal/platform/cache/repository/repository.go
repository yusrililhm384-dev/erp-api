package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type Repository struct {
	client *redis.Client
}

func (r *Repository) SetJSON(ctx context.Context, key string, val any, ttl time.Duration) error {
	body, err := json.Marshal(val)

	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, body, ttl).Err()
}

func (r *Repository) GetJSON(ctx context.Context, key string, dest any) error {
	res, err := r.client.Get(ctx, key).Result()

	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(res), &dest)
}

func (r *Repository) Set(ctx context.Context, key string, val string, ttl time.Duration) error {
	return r.client.Set(ctx, key, val, ttl).Err()
}

func (r *Repository) IncrExpired(ctx context.Context, key string, ttl time.Duration) (int64, error) {
	count, err := r.client.Incr(ctx, key).Result()

	if err != nil {
		return 0, err
	}

	if count == 1 {
		if err := r.client.Expire(ctx, key, ttl).Err(); err != nil {
			return 0, err
		}
	}

	return count, nil
}

func (r *Repository) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *Repository) Delete(ctx context.Context, keys ...string) error {
	return r.client.Del(ctx, keys...).Err()
}

func (r *Repository) SetNX(ctx context.Context, key string, val any, ttl time.Duration) error {
	return r.client.SetNX(ctx, key, val, ttl).Err()
}

func (r *Repository) Incr(ctx context.Context, key string) error {
	return r.client.Incr(ctx, key).Err()
}

func New(client *redis.Client) *Repository {
	return &Repository{
		client: client,
	}
}
