package cache

import (
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	host    string
	db      int
	Expires time.Duration
}

func NewRedisCache(host string, db int, exp time.Duration) *RedisCache {
	return &RedisCache{
		host:    host,
		db:      db,
		Expires: exp,
	}
}

func (cache *RedisCache) GetClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cache.host,
		Password: "",
		DB:       cache.db,
	})
}
