package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisStore struct {
	rdb *redis.Client
}

func NewRedisStore(rdb *redis.Client) *redisStore {
	return &redisStore{
		rdb: rdb,
	}
}

func (r *redisStore) SaveResetToken(ctx context.Context, token, userID string) error {
	return r.rdb.Set(ctx, "reset:"+token, userID, 15*time.Minute).Err()
}

func (r *redisStore) GetResetToken(ctx context.Context, token string) (string, error) {
	return r.rdb.Get(ctx, "reset:"+token).Result()
}

func (r *redisStore) DeleteResetToken(ctx context.Context, token string) error {
	return r.rdb.Del(ctx, "reset:"+token).Err()
}