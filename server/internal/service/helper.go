package service

import (
	"crypto/rand"
	"encoding/base64"
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
