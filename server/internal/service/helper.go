package service

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func generateAPIKey() (string, error) {
	b := make([]byte, 32) // 32 bytes = 256 bits

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return "sk_" + hex.EncodeToString(b), nil
}

func generateToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func VerifyPassword(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func isDuplicateError(err error) bool {
	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code == "23505"
	}
	return false
}
