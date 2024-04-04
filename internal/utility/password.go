package utility

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func ValidatePassword(hash []byte, password []byte) error {
	return bcrypt.CompareHashAndPassword(hash, password)
}
