package middleware

import (
	"context"
	"net/http"
	"strings"

	. "github.com/bernardinorafael/internal/_shared/errors"
	"github.com/bernardinorafael/internal/infra/token"
	"github.com/bernardinorafael/pkg/logger"
)

type AuthKey struct{}

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

func (m *middleware) WithAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken := r.Header.Get("Authorization")

		if len(accessToken) == 0 {
			NewHttpError(w, NewUnauthorizedError("access token not provided", nil))
			return
		}

		claims, err := token.Verify(m.secretKey, accessToken)
		if err != nil {
			if strings.Contains(err.Error(), "token has expired") {
				NewHttpError(w, NewUnauthorizedError("token has expired", err))
				return
			}
			NewHttpError(w, NewUnauthorizedError("invalid access token", err))
			return
		}
		ctx := context.WithValue(r.Context(), AuthKey{}, claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
