package usecases

import (
	"fmt"
	"spotigram/internal/customerrors"
	"spotigram/internal/service/abstractions"
	"spotigram/internal/service/models"
	"spotigram/internal/utility"
)

// A use case for changing a user's name.
// Expects access token deserialization beforehand.
// Validates the passed uuid and name.
// May return ErrInvalidInput, ErrInternal, ErrNotFound on failure.
func ChangeName(cni models.ChangeNameInput) error {
	if check := utility.IsValidUUID(cni.UserUUID); !check {
		return &customerrors.ErrInvalidInput{
			Message: "invalid \"uuid\""}
	}

	valid := utility.IsValidStructField(cni, "Name")
	if !valid {
		return &customerrors.ErrInvalidInput{
			Message: "invalid \"name\" (must be 8-100 chars long)"}
	}

	err := abstractions.UserRepositoryInstance.UpdateName(
		cni.UserUUID, cni.Name)

	return err
}

// A use case for changing a user's password.
// Expects access token deserialization beforehand.
// Validates the passed uuid and passwords.
// May return ErrInvalidInput, ErrInternal, ErrNotFound on failure.
func ChangePassword(cpi models.ChangePasswordInput) error {
	if check := utility.IsValidUUID(cpi.UserUUID); !check {
		return &customerrors.ErrInvalidInput{
			Message: "invalid \"uuid\""}
	}

	valid := utility.IsValidStructField(cpi, "OldPassword")
	if !valid {
		return &customerrors.ErrInvalidInput{
			Message: "invalid \"old_password\" (must be 8-72 chars long)"}
	}
	valid = utility.IsValidStructField(cpi, "NewPassword")
	if !valid {
		return &customerrors.ErrInvalidInput{
			Message: "invalid \"new_password\" (must be 8-72 chars long)"}
	}
	valid = utility.IsValidStructField(cpi, "NewPasswordConfirmed")
	if !valid {
		return &customerrors.ErrInvalidInput{
			Message: "invalid \"new_password_confirmed\" (must be 8-72 chars long)"}
	}

	oldPasswordHash, err := abstractions.UserRepositoryInstance.GetPassword(cpi.UserUUID)
	if err != nil {
		return err
	}

	err = utility.ValidatePassword([]byte(oldPasswordHash), []byte(cpi.OldPassword))
	if err != nil {
		return &customerrors.ErrInvalidInput{
			Message: "incorrect old password"}
	}

	if cpi.NewPassword != cpi.NewPasswordConfirmed {
		return &customerrors.ErrInvalidInput{
			Message: "new passwords do not match up"}
	}

	newPasswordHash, err := utility.HashPassword(cpi.NewPassword)
	if err != nil {
		panic(fmt.Errorf("fatal error hashing sign up password: %v", err))
	}

	err = abstractions.UserRepositoryInstance.UpdatePassword(
		cpi.UserUUID, newPasswordHash)

	return err
}

// A use case for changing a user's public key (base64).
// Expects access token deserialization beforehand.
// Validates the passed uuid and public key.
// May return ErrInvalidInput, ErrInternal, ErrNotFound on failure.
func ChangePublicKey(cpki models.ChangePublicKeyInput) error {
	if check := utility.IsValidUUID(cpki.UserUUID); !check {
		return &customerrors.ErrInvalidInput{
			Message: "invalid \"uuid\""}
	}

	valid := utility.IsValidStructField(cpki, "PublicKey")
	if cpki.PublicKey == "" || !valid {
		return &customerrors.ErrInvalidInput{
			Message: "invalid \"public_key\" (must be 1-6120 chars long)"}
	}

	err := abstractions.UserRepositoryInstance.UpdatePublicKey(
		cpki.UserUUID, cpki.PublicKey)

	return err
}

// A use case for changing a user's picture (raw bytes).
// Expects access token deserialization beforehand.
// Validates the passed uuid and image_array.
// May return ErrInvalidInput, ErrInternal, ErrNotFound on failure.
func ChangePicture(cpi models.ChangePictureInput) error {
	if check := utility.IsValidUUID(cpi.UserUUID); !check {
		return &customerrors.ErrInvalidInput{
			Message: "invalid \"uuid\""}
	}

	if cpi.Image == nil || len(cpi.Image) == 0 || len(cpi.Image) > 5*1024*1024 {
		return &customerrors.ErrInvalidInput{
			Message: "invalid image, must be a png or a jpg (under 5 megabytes)"}
	}

	imageWebP, err := utility.ConvertAndResizeImageToWebP(cpi.Image, 512, 512)
	if err != nil {
		return &customerrors.ErrInvalidInput{
			Message: "invalid image, must be a png or a jpg (under 5 megabytes)"}
	}

	err = abstractions.UserRepositoryInstance.UpdatePicture(
		cpi.UserUUID, imageWebP)

	return err
}
