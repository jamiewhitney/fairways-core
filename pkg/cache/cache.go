package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type Config struct {
	RedisHost     string `env:"REDIS_HOST, default=localhost" json:"redis_host,omitempty"`
	RedisPort     string `env:"REDIS_PORT, default=6379" json:"redis_port,omitempty"`
	RedisUsername string `env:"REDIS_USERNAME" json:"redis_username,omitempty"`
	RedisPassword string `env:"REDIS_PASSWORD" json:"-,omitempty"`
}

type RedisRepository struct {
	*redis.Client
	expiry time.Duration
}

func New(address string, username, password string, expiry time.Duration) (*RedisRepository, error) {
	Client := redis.NewClient(&redis.Options{
		Addr:      address,
		TLSConfig: nil,
	})

	if err := Client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return &RedisRepository{Client, expiry}, nil
}
