package usecases

import (
	"spotigram/internal/customerrors"
	"spotigram/internal/service/abstractions"
	"spotigram/internal/service/models"
	"spotigram/internal/utility"
)

// A use case to get a user list.
// Expects access token deserialization beforehand.
// May return ErrInvalidInput, ErrInternal, ErrNotFound on failure.
func GetUsers(input models.GetUsersInput) ([]models.User, error) {
	if check := input.Offset >= 0; !check {
		return nil, &customerrors.ErrInvalidInput{
			Message: "invalid \"offset\""}
	}
	users, err :=
		abstractions.UserRepositoryInstance.GetUsers(input.Offset)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// A use case for current user info.
// Expects access token deserialization beforehand.
// Validates the passed uuid.
// May return ErrInvalidInput, ErrInternal, ErrNotFound on failure.
func GetUserInfo(mii models.GetUserInfoInput) (*models.User, error) {
	if check := utility.IsValidUUID(mii.UserUUID); !check {
		return nil, &customerrors.ErrInvalidInput{
			Message: "invalid \"uuid\""}
	}
	user, err :=
		abstractions.UserRepositoryInstance.GetUser(mii.UserUUID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// A use case for current user's public key.
// Expects access token deserialization beforehand.
// Validates the passed uuid.
// May return ErrInvalidInput, ErrInternal, ErrNotFound on failure.
func GetPublicKey(mii models.GetPublicKeyInput) (string, error) {
	if check := utility.IsValidUUID(mii.UserUUID); !check {
		return "", &customerrors.ErrInvalidInput{
			Message: "invalid \"uuid\""}
	}
	key, err :=
		abstractions.UserRepositoryInstance.GetPublicKey(mii.UserUUID)
	if err != nil {
		return "", err
	}
	return key, nil
}

// A use case for current user's public key.
// Expects access token deserialization beforehand.
// Validates the passed uuid.
// May return ErrInvalidInput, ErrInternal, ErrNotFound on failure.
func GetPicture(mii models.GetPictureInput) ([]byte, error) {
	if check := utility.IsValidUUID(mii.UserUUID); !check {
		return nil, &customerrors.ErrInvalidInput{
			Message: "invalid \"uuid\""}
	}
	pic, err :=
		abstractions.UserRepositoryInstance.GetPicture(mii.UserUUID)
	if err != nil {
		return nil, err
	}
	return pic, nil
}
