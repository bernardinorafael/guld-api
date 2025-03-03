package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"

	"runtime/debug"

	"github.com/bernardinorafael/internal/_shared/errors"
)

func WithRecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("panic recovered: %+v\nStack trace: %s\n", err, debug.Stack())

				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("Connection", "close")

				appErr := errors.NewInternalServerError(fmt.Errorf("internal server error: %v", err))
				w.WriteHeader(appErr.StatusCode())

				if encodeErr := json.NewEncoder(w).Encode(appErr); encodeErr != nil {
					fmt.Printf("error encoding response: %v\n", encodeErr)
				}
			}
		}()
		next.ServeHTTP(w, r)
	})
}
