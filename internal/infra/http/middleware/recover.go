package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bernardinorafael/internal/_shared/errors"
)

func WithRecoverPanic(done http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				appErr := errors.NewInternalServerError(fmt.Errorf("%s", err))

				w.Header().Set("Connection", "close")
				w.WriteHeader(appErr.StatusCode())
				_ = json.NewEncoder(w).Encode(appErr)
			}
		}()
		done.ServeHTTP(w, r)
	})
}
