package token

import (
	"context"
	"errors"
	"strings"

	"github.com/aead/chacha20poly1305"
	. "github.com/bernardinorafael/internal/_shared/errors"
	"github.com/bernardinorafael/pkg/logger"
	"github.com/golang-jwt/jwt"
)

type Token struct {
	ctx       context.Context
	log       logger.Logger
	secretKey string
}

func New(ctx context.Context, log logger.Logger, secretKey string) *Token {
	return &Token{ctx, log, secretKey}
}

func (t *Token) GenToken(params WithParams) (string, *AccountClaims, error) {
	var token string

	if len(t.secretKey) != chacha20poly1305.KeySize {
		msg := "invalid secret key"
		t.log.Errorw(t.ctx, msg, logger.Err(errors.New(msg)))
		return token, nil, NewBadRequestError(msg, nil)
	}

	claims, err := NewAccountClaims(params)
	if err != nil {
		msg := "failed to create account claims"
		t.log.Errorw(t.ctx, msg, logger.Err(err))
		return token, nil, NewBadRequestError(msg, err)
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = jwtToken.SignedString([]byte(t.secretKey))
	if err != nil {
		msg := "failed to sign token"
		t.log.Errorw(t.ctx, msg, logger.Err(err))
		return token, claims, NewBadRequestError(msg, err)
	}

	return token, claims, nil
}

func (t *Token) VerifyToken(v string) (*AccountClaims, error) {
	// This is a sanity check to ensure the token is not empty
	if strings.TrimSpace(v) == "" {
		msg := "invalid token"
		t.log.Errorw(t.ctx, msg, logger.Err(errors.New(msg)))
		return nil, NewBadRequestError(msg, nil)
	}

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			msg := "invalid token signing method"
			t.log.Errorw(t.ctx, msg, logger.Err(errors.New(msg)))
			return nil, NewUnauthorizedError(msg, nil)
		}
		return []byte(t.secretKey), nil
	}

	token, err := jwt.ParseWithClaims(v, &AccountClaims{}, keyFunc)
	if err != nil {
		msg := "failed to parse token"
		t.log.Errorw(t.ctx, msg, logger.Err(err))
		return nil, NewBadRequestError(msg, err)
	}

	claims, ok := token.Claims.(*AccountClaims)
	if !ok {
		msg := "invalid token claims"
		t.log.Errorw(t.ctx, msg, logger.Err(errors.New(msg)))
		return nil, NewBadRequestError(msg, nil)
	}

	return claims, nil
}
