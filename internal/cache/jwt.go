package cache

import (
	"context"
	"time"

	"spotigram/internal/customerrors"

	"github.com/redis/go-redis/v9"
)

type JWTRedisCache struct {
	RedisClient *redis.Client
}

// Deletes refresh and access tokens (order does not matter).
// Returns the number of deleted tokens.
// May return ErrInternal on failure.
func (j *JWTRedisCache) DeleteRefreshAndAccessToken(
	refreshTokenUUID, accessTokenUUID string) (int64, error) {
	ctx := context.TODO()
	num, err := j.RedisClient.Del(ctx, refreshTokenUUID, accessTokenUUID).Result()
	if err != nil {
		return 0, &customerrors.ErrInternal{Message: err.Error()}
	}
	return num, nil
}

// Returns value of the key, where key is meant to be
// a uuid of a JWT, value is meant to be user uuid.
// May return ErrNotFound or ErrInternal on failure.
func (j *JWTRedisCache) GetToken(uuid string) (string, error) {
	ctx := context.TODO()
	userUUID, err := j.RedisClient.Get(ctx, uuid).Result()
	if err != nil {
		if err == redis.Nil {
			return "", &customerrors.ErrNotFound{Message: err.Error()}
		}
		return "", &customerrors.ErrInternal{Message: err.Error()}
	}
	return userUUID, nil
}

// Sets key-value pair with an expiration time.
// Key is meant to be a uuid of a JWT, value is meant to be user uuid.
// May return ErrInternal on failure.
func (j *JWTRedisCache) SetToken(
	key string, value string, expiresIn time.Duration) error {
	ctx := context.TODO()
	err := j.RedisClient.Set(ctx,
		key, value, expiresIn).Err()
	if err != nil {
		return &customerrors.ErrInternal{Message: err.Error()}
	}
	return nil
}
