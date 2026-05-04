package redis

import (
	"fmt"
	"strconv"

	"github.com/diyorbek/sentinel/internal/config"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(cfg *config.Config) *redis.Client {
	addr := fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort)
	db, _ := strconv.Atoi(cfg.RedisDB)

	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.RedisPassword,
		DB:       db, 
	})
}
