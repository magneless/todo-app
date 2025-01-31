package hashing

import (
	"crypto/sha1"
	"fmt"
)

const salt = "afasfafmvpmvoermbpqoa123"

func HashPassword(password string) (string, error) {
	hash := sha1.New()
	_, err := hash.Write([]byte(password + salt))
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return fmt.Sprintf("%x", hash.Sum([]byte(salt))), err
}

func VerifyPassword(password, passwordHash string) (bool, error) {
	hash, err := HashPassword(password)
	if err != nil {
		return false, fmt.Errorf("failed to verify password: %w", err)
	}
	return passwordHash == hash, err
}
