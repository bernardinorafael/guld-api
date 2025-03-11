package token

import (
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt"
)

func Verify(secretKey string, v string) (*AccountClaims, error) {
	if strings.TrimSpace(v) == "" {
		return nil, fmt.Errorf("invalid token")
	}

	keyFunc := func(token *jwt.Token) (any, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("invalid token signing method")
		}
		return []byte(secretKey), nil
	}

	token, err := jwt.ParseWithClaims(v, &AccountClaims{}, keyFunc)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*AccountClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
