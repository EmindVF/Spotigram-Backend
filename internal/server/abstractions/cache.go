package abstractions

import "time"

type JWTCache interface {
	DeleteRefreshAndAccessToken(string, string) (int64, error)
	GetToken(string) (string, error)
	SetToken(key string, value string, expiresIn time.Duration) error
}
