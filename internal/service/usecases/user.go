package usecases

import (
	"fmt"
	"strings"

	"spotigram/internal/customerrors"
	"spotigram/internal/service/abstractions"
	"spotigram/internal/service/models"
	"spotigram/internal/utility"
)

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

func SignInUser(sii models.SignInInput) (uuid string, e error) {
	sii.Email = strings.TrimSpace(sii.Email)

	valid := utility.IsValidStructField(sii, "Email")
	if !valid {
		return "", &customerrors.ErrInvalidInput{
			Message: "invalid \"email\" (must be 5-100 chars long)"}
	}

	valid = utility.IsValidStructField(sii, "Password")
	if !valid {
		return "", &customerrors.ErrInvalidInput{
			Message: "invalid \"password\" (must be 8-72 chars long)"}
	}

	uuid, passwordHash, err := abstractions.UserRepositoryInstance.GetUUIDAndPasswordByEmail(sii.Email)
	if err != nil {
		return "", err
	}

	err = utility.ValidatePassword([]byte(passwordHash), []byte(sii.Password))
	if err != nil {
		return "", &customerrors.ErrInvalidInput{
			Message: "incorrect password"}
	}

	return uuid, nil
}

func DoesUserExist(uuid string) (bool, error) {
	if check := utility.IsValidUUID(uuid); !check {
		return false, &customerrors.ErrInvalidInput{
			Message: "invalid \"uuid\""}
	}
	exists, err := abstractions.UserRepositoryInstance.DoesUserExist(uuid)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func GetUser(uuid string) (*models.User, error) {
	if check := utility.IsValidUUID(uuid); !check {
		return nil, &customerrors.ErrInvalidInput{
			Message: "invalid \"uuid\""}
	}
	user, err := abstractions.UserRepositoryInstance.GetUser(uuid)
	if err != nil {
		return nil, err
	}
	return user, nil
}
