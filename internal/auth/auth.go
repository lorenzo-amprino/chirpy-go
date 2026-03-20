package auth

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/alexedwards/argon2id"
)

func HashPassword(password string) (string, error) {
	hashed, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return hashed, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, err
	}
	return match, nil
}

func GetAPIKey(headers http.Header) (string, error) {
	apiKey := headers.Get("Authorization")
	if apiKey == "" {
		return "", errors.New("missing Authorization header")
	}
	var token string
	fmt.Sscanf(apiKey, "ApiKey %s", &token)
	if token == "" {
		return "", errors.New("invalid Authorization header format")
	}
	return token, nil
}
