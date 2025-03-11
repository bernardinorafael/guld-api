package middleware

import (
	"context"
	"net/http"

	. "github.com/bernardinorafael/internal/_shared/errors"
	"github.com/bernardinorafael/internal/infra/token"
	"github.com/bernardinorafael/pkg/logger"
)

type Key string

const (
	UserIDKey Key = "user_id"
	AccIDKey  Key = "acc_id"
)

type middleware struct {
	log       logger.Logger
	secretKey string
}

func NewWithAuth(log logger.Logger, secretKey string) *middleware {
	return &middleware{
		log:       log,
		secretKey: secretKey,
	}
}

func (m *middleware) WithAuth(done http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken := r.Header.Get("Authorization")

		if len(accessToken) == 0 {
			NewHttpError(w, NewUnauthorizedError("access token not provided", nil))
			return
		}

		p, err := token.Verify(m.secretKey, accessToken)
		if err != nil {
			NewHttpError(w, NewUnauthorizedError("invalid access token", err))
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, p.UserID)
		ctx = context.WithValue(ctx, AccIDKey, p.AccountID)

		done.ServeHTTP(w, r.WithContext(ctx))
	})
}
