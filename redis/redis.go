package redis

import "github.com/go-redis/redis/v8"

var (
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
)
