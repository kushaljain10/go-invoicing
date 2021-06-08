package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/kushaljain/go-invoicing/taxes"
)

type redisCache struct {
	host    string
	db      int
	expires time.Duration
}

func NewRedisCache(host string, db int, exp time.Duration) Cache {
	return &redisCache{
		host:    host,
		db:      db,
		expires: exp,
	}
}

func (cache *redisCache) getClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cache.host,
		Password: "",
		DB:       cache.db,
	})
}

func (cache *redisCache) Set(key string, value *taxes.Taxes) {
	ctx := context.TODO()
	client := cache.getClient()

	json, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}

	client.Set(ctx, key, json, cache.expires*time.Second)
}

func (cache *redisCache) Get(key string) *taxes.Taxes {
	ctx := context.TODO()
	client := cache.getClient()

	val, err := client.Get(ctx, key).Result()
	if err != nil {
		return nil
	}

	tax := taxes.Taxes{}
	err = json.Unmarshal([]byte(val), &tax)
	if err != nil {
		panic(err)
	}
	return &tax
}
