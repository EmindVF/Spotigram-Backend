package abstractions

import "spotigram/internal/service/models"

type UserRepository interface {
	AddUser(models.User) error
	GetUUIDAndPasswordByEmail(string) (string, string, error)
	DoesUserExist(string) (bool, error)
	GetUser(string) (*models.User, error)
}
