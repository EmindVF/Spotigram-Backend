package utility

import "golang.org/x/crypto/bcrypt"

// Hashes a password using bcrypt.
// Password's max length is 72.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// Compares password to its bcrypt hash.
// Returns nil on validation, error otherwise.
func ValidatePassword(hash []byte, password []byte) error {
	return bcrypt.CompareHashAndPassword(hash, password)
}
