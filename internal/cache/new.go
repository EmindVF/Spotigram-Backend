package cache

import (
	"spotigram/internal/service/abstractions"
)

// Returns a redis based JWT cache.
func NewJWTCache() abstractions.JWTCache {
	return &JWTRedisCache{RedisClient: RedisClient}
}
