package loggerconf

import (
	"context"

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
	return args
}
