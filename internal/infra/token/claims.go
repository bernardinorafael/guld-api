package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/segmentio/ksuid"
)

type AccountClaims struct {
	AccountID string `json:"account_id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	jwt.RegisteredClaims
}

func NewAccountClaims(accountId, email, username string, duration time.Duration) (*AccountClaims, error) {
	claims := &AccountClaims{
		AccountID: accountId,
		Username:  username,
		Email:     email,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        ksuid.New().String(),
			Subject:   email,
			Audience:  []string{},
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}

	return claims, nil
}

func (a *AccountClaims) Valid() error {
	if a.ExpiresAt != nil && !a.ExpiresAt.After(time.Now()) {
		return errors.New("token has expired")
	}
	return nil
}
