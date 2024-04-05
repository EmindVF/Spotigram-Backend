package usecases

import (
	"spotigram/internal/customerrors"
	"spotigram/internal/service/abstractions"
	"spotigram/internal/service/models"
	"spotigram/internal/utility"
)

// A use case for current user info.
// Validates the passed uuid.
// May return ErrInvalidInput, ErrInternal, ErrNotFound on failure.
func MyInfo(mii models.MyInfoInput) (*models.User, error) {
	if check := utility.IsValidUUID(mii.AccesssTokenUUID); !check {
		return nil, &customerrors.ErrInvalidInput{
			Message: "invalid \"uuid\""}
	}
	user, err :=
		abstractions.UserRepositoryInstance.GetUser(mii.AccesssTokenUUID)
	if err != nil {
		return nil, err
	}
	return user, nil
}
