package abstractions

import "spotigram/internal/service/models"

type UserRepository interface {
	// Adds user to the repository.
	// May return ErrInternal or ErrInvalidInput on failure.
	AddUser(models.User) error

	// Returns uuid and hashed password of an user by its email.
	// Email validation is not provided.
	// May return ErrInternal or ErrNotFound on failure.
	GetUUIDAndPasswordByEmail(string) (string, string, error)

	// Returns bool on whether the user uuid is present.
	// UUID validation is not provided.
	// May return ErrInternal on failure.
	DoesUserExist(string) (bool, error)

	// Returns a user by its uuid.
	// UUID validation is not provided.
	// May return ErrInternal or ErrNotFound on failure.
	GetUser(string) (*models.User, error)
}
