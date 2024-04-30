package usecases

import (
	"fmt"
	"strings"
	"time"

	"spotigram/internal/config"
	"spotigram/internal/customerrors"
	"spotigram/internal/service/abstractions"
	"spotigram/internal/service/models"
	"spotigram/internal/utility"
)

// A use case for user sign up.
// Provides input validation.
// May return ErrInvalidInput or ErrInternal on failure.
func SignUpUser(sui models.SignUpInput) (e error) {
	sui.Name = strings.TrimSpace(sui.Name)
	sui.Email = strings.TrimSpace(sui.Email)

	valid := utility.IsValidStructField(sui, "Name")
	if !valid {
		return &customerrors.ErrInvalidInput{
			Message: "invalid \"name\" (must be 8-100 chars long)"}
	}

	valid = utility.IsValidStructField(sui, "Email")
	if !valid {
		return &customerrors.ErrInvalidInput{
			Message: "invalid \"email\" (must be 5-100 chars long)"}
	}

	valid = utility.IsValidStructField(sui, "Password")
	if !valid {
		return &customerrors.ErrInvalidInput{
			Message: "invalid \"password\" (must be 8-72 chars long)"}
	}

	valid = utility.IsValidStructField(sui, "PasswordConfirmed")
	if !valid {
		return &customerrors.ErrInvalidInput{
			Message: "invalid \"password_confirmed\" (must be 8-72 chars long)"}
	}

	if sui.Password != sui.PasswordConfirmed {
		return &customerrors.ErrInvalidInput{
			Message: "passwords don't match up"}
	}

	hashedPassword, err := utility.HashPassword(sui.Password)
	if err != nil {
		panic(fmt.Errorf("fatal error hashing sign up password: %v", err))
	}

	err = abstractions.UserRepositoryInstance.AddUser(models.User{
		Id:       utility.GenerateUUID(),
		Name:     sui.Name,
		Email:    sui.Email,
		Password: hashedPassword,
		Verified: false,
	})

	if err != nil {
		return err
	}

	return nil
}

// A use case for user sign in.
// Returns user uuid and newly created access and refresh token details.
// Provides input validation.
// May return ErrInvalidInput, ErrInternal, ErrNotFound on failure.
func SignInUser(sii models.SignInInput, cfg *config.Config) (
	uuid string, accessTD, refreshTD *utility.TokenDetails, err error) {
	sii.Email = strings.TrimSpace(sii.Email)

	valid := utility.IsValidStructField(sii, "Email")
	if !valid {
		return "", nil, nil, &customerrors.ErrInvalidInput{
			Message: "invalid \"email\" (must be 5-100 chars long)"}
	}

	valid = utility.IsValidStructField(sii, "Password")
	if !valid {
		return "", nil, nil, &customerrors.ErrInvalidInput{
			Message: "invalid \"password\" (must be 8-72 chars long)"}
	}

	uuid, passwordHash, err :=
		abstractions.UserRepositoryInstance.GetUUIDAndPasswordByEmail(sii.Email)
	if err != nil {
		return "", nil, nil, err
	}

	err = utility.ValidatePassword([]byte(passwordHash), []byte(sii.Password))
	if err != nil {
		return "", nil, nil, &customerrors.ErrInvalidInput{
			Message: "incorrect password"}
	}

	accessTD, err = utility.CreateToken(
		uuid, cfg.AccessToken.ExpiresIn, cfg.AccessToken.PrivateKey)
	if err != nil {
		return "", nil, nil, &customerrors.ErrInternal{Message: err.Error()}
	}

	refreshTD, err = utility.CreateToken(
		uuid, cfg.RefreshToken.ExpiresIn, cfg.RefreshToken.PrivateKey)
	if err != nil {
		return "", nil, nil, &customerrors.ErrInternal{Message: err.Error()}
	}

	now := time.Now()

	err = abstractions.JWTCacheInstance.SetToken(
		accessTD.TokenUUID,
		uuid,
		(time.Unix(accessTD.ExpiresIn, 0).Sub(now)))
	if err != nil {
		return "", nil, nil, &customerrors.ErrInternal{Message: err.Error()}
	}

	err = abstractions.JWTCacheInstance.SetToken(
		refreshTD.TokenUUID,
		uuid,
		(time.Unix(refreshTD.ExpiresIn, 0).Sub(now)))
	if err != nil {
		return "", nil, nil, &customerrors.ErrInternal{Message: err.Error()}
	}

	return uuid, accessTD, refreshTD, nil
}

// A use case for user logout.
// Expects access token deserialization beforehand.
// Validates refresh token, deletes passed access and refresh tokens.
// May return ErrUnauthorized, ErrInternal on failure.
func Logout(loi models.LogoutInput, cfg *config.Config) error {
	if loi.RefreshToken == "" {
		return &customerrors.ErrUnauthorized{
			Message: "refresh token is invalid or session has expired"}
	}

	refreshTokenClaims, err := utility.ValidateToken(
		loi.RefreshToken, cfg.RefreshToken.PublicKey)
	if err != nil {
		return &customerrors.ErrUnauthorized{
			Message: err.Error()}
	}

	_, err = abstractions.JWTCacheInstance.DeleteRefreshAndAccessToken(
		refreshTokenClaims.TokenUUID, loi.AccessTokenUUID)
	if err != nil {
		return &customerrors.ErrInternal{
			Message: err.Error()}
	}
	return nil
}

// A use case for access token refresh.
// Validates refresh token.
// May return ErrUnauthorized, ErrInternal, ErrNotFound on failure.
func RefreshAccessToken(rati models.RefreshAccessTokenInput, cfg *config.Config) (
	*utility.TokenDetails, error) {

	if rati.RefreshToken == "" {
		return nil, &customerrors.ErrUnauthorized{
			Message: "could not refresh access token"}
	}

	tokenClaims, err := utility.ValidateToken(
		rati.RefreshToken, cfg.RefreshToken.PublicKey)
	if err != nil {
		return nil, &customerrors.ErrUnauthorized{
			Message: "could not refresh access token"}
	}

	userUuid, err :=
		abstractions.JWTCacheInstance.GetToken(tokenClaims.TokenUUID)
	if err != nil {
		return nil, err
	}

	if check := utility.IsValidUUID(userUuid); !check {
		return nil, &customerrors.ErrNotFound{
			Message: "invalid \"uuid\""}
	}
	exists, err :=
		abstractions.UserRepositoryInstance.DoesUserExist(userUuid)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, &customerrors.ErrNotFound{
			Message: "user belonging to this token no longer exists"}
	}

	accessTokenDetails, err := utility.CreateToken(
		userUuid, cfg.AccessToken.ExpiresIn, cfg.AccessToken.PrivateKey)
	if err != nil {
		return nil, &customerrors.ErrInternal{Message: err.Error()}
	}

	now := time.Now()

	err = abstractions.JWTCacheInstance.SetToken(
		accessTokenDetails.TokenUUID,
		userUuid,
		(time.Unix(accessTokenDetails.ExpiresIn, 0).Sub(now)))
	if err != nil {
		return nil, &customerrors.ErrInternal{Message: err.Error()}
	}

	return accessTokenDetails, nil
}
