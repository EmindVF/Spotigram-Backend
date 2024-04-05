package cache

import (
	"context"
	"fmt"

	"spotigram/internal/config"

	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *redis.Client
)

// Connects package's redis client to the redis
// server from the url in the config file.
func ConnectRedis(cfg *config.Config) {
	ctx := context.TODO()

	RedisClient = redis.NewClient(&redis.Options{
		Addr: cfg.Cache.RedisUrl,
	})

	if _, err := RedisClient.Ping(ctx).Result(); err != nil {
		panic(fmt.Errorf("error pinging redis server: %v", err))
	}
}
