package loggerconf

import (
	"context"

	"github.com/bernardinorafael/internal/infra/http/middleware"
	"github.com/bernardinorafael/pkg/logger"
)

func New(app string, debugLevel bool) logger.Logger {
	params := logger.LogParams{
		AppName:                  app,
		DebugLevel:               debugLevel,
		AddAttributesFromContext: addDefaultAttrToLogger,
	}
	return logger.New(params)
}

func addDefaultAttrToLogger(ctx context.Context) []logger.LogField {
	args := []logger.LogField{}

	if userID, ok := getUserID(ctx); ok {
		args = append(args, logger.String(string(middleware.UserIDKey), userID))
	}

	return args
}

func getUserID(ctx context.Context) (string, bool) {
	if ctx == nil {
		return "", false
	}

	v := ctx.Value(middleware.UserIDKey)
	if v == nil {
		return "", false
	}

	return v.(string), false
}
