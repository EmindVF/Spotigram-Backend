package utility

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Details of an JWT.
type TokenDetails struct {
	Token     string
	TokenUUID string
	UserUUID  string
	ExpiresIn int64
}

// Creates a JWT, returning its details.
func CreateToken(uuid string, ttl time.Duration, privateKey []byte) (*TokenDetails, error) {
	now := time.Now().UTC()
	td := &TokenDetails{
		ExpiresIn: now.Add(ttl).Unix(),
		TokenUUID: GenerateUUID(),
		UserUUID:  uuid,
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		panic(fmt.Errorf("error parsing private key to jwt: %v", err))
	}

	atClaims := make(jwt.MapClaims)
	atClaims["sub"] = uuid
	atClaims["token_uuid"] = td.TokenUUID
	atClaims["exp"] = td.ExpiresIn
	atClaims["iat"] = now.Unix()
	atClaims["nbf"] = now.Unix()

	td.Token, err = jwt.NewWithClaims(jwt.SigningMethodRS256, atClaims).SignedString(key)
	if err != nil {
		panic(fmt.Errorf("error creating jwt token: %v", err))
	}
	return td, nil
}

// Validates a JWT, returning its details, containing
// only UserUUID and TokenUUID.
func ValidateToken(token string, publicKey []byte) (*TokenDetails, error) {
	key, err := jwt.ParseRSAPublicKeyFromPEM(publicKey)
	if err != nil {
		panic(fmt.Errorf("error parsing public key to jwt: %v", err))
	}

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", t.Header["alg"])
		}
		return key, nil
	})

	if err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, fmt.Errorf("validation error: invalid token")
	}

	return &TokenDetails{
		UserUUID:  fmt.Sprint(claims["sub"]),
		TokenUUID: fmt.Sprint(claims["token_uuid"]),
	}, nil
}
