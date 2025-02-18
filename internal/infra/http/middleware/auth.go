package middleware

import (
	"context"
	"net/http"

	"github.com/bernardinorafael/internal/_shared/errors"
	"github.com/bernardinorafael/internal/infra/token"
	"github.com/bernardinorafael/pkg/logger"
)

type Key string

const (
	UserIDKey Key = "user_id"
	AccIDKey  Key = "acc_id"
)

type middleware struct {
	ctx context.Context
	log logger.Logger
	t   *token.Token
}

func NewWithAuth(ctx context.Context, log logger.Logger, secretKey string) *middleware {
	return &middleware{
		ctx: ctx,
		log: log,
		t:   token.New(ctx, log, secretKey),
	}
}

func (m *middleware) WithAuth(done http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken := r.Header.Get("Authorization")

		if len(accessToken) == 0 {
			errors.NewHttpError(w, errors.NewUnauthorizedError("access token not provided", nil))
			return
		}

		p, err := m.t.VerifyToken(accessToken)
		if err != nil {
			errors.NewHttpError(w, errors.NewUnauthorizedError("invalid access token", err))
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, p.UserID)
		ctx = context.WithValue(ctx, AccIDKey, p.AccountID)

		done.ServeHTTP(w, r.WithContext(ctx))
	})
}
