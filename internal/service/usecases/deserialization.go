package usecases

import (
	"spotigram/internal/config"
	"spotigram/internal/customerrors"
	"spotigram/internal/service/abstractions"
	"spotigram/internal/service/models"
	"spotigram/internal/utility"
)

// A use case for access token deserialization.
// Validates access token.
// May return ErrUnauthorized, ErrInternal, ErrNotFound on failure.
func DeserializeToken(dti models.DeserializeTokenInput, cfg *config.Config) (
	userUuid string, accessTokenUuid string, err error) {

	if dti.AccessToken == "" {
		return "", "", &customerrors.ErrUnauthorized{
			Message: "you are not logged in"}
	}

	tokenClaims, err := utility.ValidateToken(
		dti.AccessToken, cfg.AccessToken.PublicKey)
	if err != nil {
		return "", "", &customerrors.ErrUnauthorized{
			Message: err.Error()}
	}

	userUuid, err = abstractions.JWTCacheInstance.GetToken(tokenClaims.TokenUUID)
	if err != nil {
		return "", "", &customerrors.ErrUnauthorized{
			Message: "token is invalid or session has expired"}
	}

	if check := utility.IsValidUUID(userUuid); !check {
		return "", "", &customerrors.ErrNotFound{
			Message: "invalid \"uuid\""}
	}
	exists, err :=
		abstractions.UserRepositoryInstance.DoesUserExist(userUuid)
	if err != nil {
		return "", "", err
	}
	if !exists {
		return "", "", &customerrors.ErrNotFound{
			Message: "user belonging to this token no longer exists"}
	}

	accessTokenUuid = tokenClaims.TokenUUID

	return userUuid, accessTokenUuid, nil
}
