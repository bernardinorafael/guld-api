package session

import (
	"context"

	"github.com/bernardinorafael/pkg/logger"
	"github.com/go-chi/chi"
)

type controller struct {
	ctx context.Context
	log logger.Logger
	// svc       ServiceInterface
	secretKey string
}

func NewController(
	ctx context.Context,
	log logger.Logger,
	// svc ServiceInterface,
	secretKey string,
) *controller {
	return &controller{
		ctx: ctx,
		log: log,
		// svc:       svc,
		secretKey: secretKey,
	}
}

func (c controller) RegisterRoute(r *chi.Mux) {}
