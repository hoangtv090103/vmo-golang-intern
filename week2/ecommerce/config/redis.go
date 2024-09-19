package config

import (
	"github.com/go-redis/redis/v8"
)

type Redis struct {
	client *redis.Client
}

func NewRedis() *Redis {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	return &Redis{
		client: rdb,
	}
}

func (r *Redis) Close() {
	r.client.Close()
}

func (r *Redis) GetClient() *redis.Client {
	return r.client
}

