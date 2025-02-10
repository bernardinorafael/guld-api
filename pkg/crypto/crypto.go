package crypto

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword encrypts a password using bcrypt
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), err
}

// PasswordMatches compares a plain text password with an encrypted password
// Returns true if the password matches, false otherwise
func PasswordMatches(password, encrypted string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(encrypted), []byte(password))
	return err == nil
}
