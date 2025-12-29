package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// IsStrongPassword checks if the password is strong
//
// A strong password is defined as one that is at least 8 characters long
// and contains at least one uppercase letter, one lowercase letter,
// one number, and one special character.
func IsStrongPassword(password string) bool {
	var hasMinLen, hasUpper, hasLower, hasNumber, hasSpecial bool
	if len(password) >= 8 {
		hasMinLen = true
	}
	for _, char := range password {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case 'a' <= char && char <= 'z':
			hasLower = true
		case '0' <= char && char <= '9':
			hasNumber = true
		case (char >= 33 && char <= 47) || (char >= 58 && char <= 64) ||
			(char >= 91 && char <= 96) || (char >= 123 && char <= 126):
			hasSpecial = true
		}
	}
	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}

// IsWeakPassword checks if the password is weak (less than 6 characters)
func IsWeakPassword(password string) bool {
	var hasMinLen bool
	if len(password) < 6 {
		return true
	}
	return !hasMinLen
}

// HashPassword hashes the given password using bcrypt
func HashPassword(password string) (string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPass), nil
}

// ComparePasswordHash compares the given password with its bcrypt hash
func ComparePasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
