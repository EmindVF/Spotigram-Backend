package abstractions

import "time"

type JWTCache interface {
	// Deletes refresh and access tokens (order does not matter).
	// Returns the number of deleted tokens.
	// May return ErrInternal on failure.
	DeleteRefreshAndAccessToken(string, string) (int64, error)

	// Returns value of the key, where key is meant to be
	// a uuid of a JWT, value is meant to be user uuid.
	// May return ErrNotFound or ErrInternal on failure.
	GetToken(string) (string, error)

	// Sets key-value pair with an expiration time.
	// Key is meant to be a uuid of a JWT, value is meant to be user uuid.
	// May return ErrInternal on failure.
	SetToken(key string, value string, expiresIn time.Duration) error
}
