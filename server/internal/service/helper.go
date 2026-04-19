package service

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

func generateAPIKey() (string, error) {
	b := make([]byte, 32) // 32 byte = kuchli key

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	// URL-safe string
	return base64.URLEncoding.EncodeToString(b), nil
}

func HassPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func VerifyPassword(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}