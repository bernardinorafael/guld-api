package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/segmentio/ksuid"
)

type WithParams struct {
	AccountID string
	UserID    string
	OrgID     *string
	Username  string
	Email     string
	Duration  time.Duration
}

type AccountClaims struct {
	AccountID string  `json:"account_id"`
	UserID    string  `json:"user_id"`
	OrgID     *string `json:"org_id"`
	Username  string  `json:"username"`
	Email     string  `json:"email"`
	jwt.RegisteredClaims
}

func NewAccountClaims(params WithParams) (*AccountClaims, error) {
	claims := &AccountClaims{
		AccountID: params.AccountID,
		UserID:    params.UserID,
		OrgID:     params.OrgID,
		Username:  params.Username,
		Email:     params.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        ksuid.New().String(),
			Subject:   params.Email,
			Audience:  []string{},
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(params.Duration)),
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
