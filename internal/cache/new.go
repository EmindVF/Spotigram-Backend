package cache

import (
	"spotigram/internal/server/abstractions"
)

func NewJWTCache() abstractions.JWTCache {
	return &JWTRedisCache{RedisClient: RedisClient}
}
