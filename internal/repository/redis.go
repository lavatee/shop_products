package repository

import (
	"time"

	"github.com/redis/go-redis/v9"
)

const cachedItemTTL = 10 * time.Hour

func NewRedisDB(host string) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: "",
		DB:       0,
	})
	return rdb
}
