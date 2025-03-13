package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AccountClaims struct {
	AccountID string `json:"account_id"`
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	jwt.RegisteredClaims
}

func NewAccountClaims(accId, userId, username string, duration time.Duration) (*AccountClaims, error) {
	claims := &AccountClaims{
		AccountID: accId,
		UserID:    userId,
		Username:  username,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   username,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}

	return claims, nil
}

func (a *AccountClaims) Valid() error {
	if time.Now().After(a.ExpiresAt.Time) {
		return errors.New("token has expired")
	}
	return nil
}
