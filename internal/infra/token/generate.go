package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/golang-jwt/jwt"
)

func Generate(secretKey, accId, userId, username string, duration time.Duration) (string, *AccountClaims, error) {
	if len(secretKey) != chacha20poly1305.KeySize {
		return "", nil, fmt.Errorf("invalid secret key")
	}

	claims, err := NewAccountClaims(accId, userId, username, duration)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create account claims: %w", err)
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := jwtToken.SignedString([]byte(secretKey))
	if err != nil {
		return "", claims, fmt.Errorf("failed to sign token: %w", err)
	}

	return token, claims, nil
}
