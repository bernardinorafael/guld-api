package org

import (
	"context"
	"net/http"

	. "github.com/bernardinorafael/internal/_shared/errors"
	"github.com/bernardinorafael/internal/_shared/util"
	"github.com/bernardinorafael/internal/infra/http/middleware"
	"github.com/bernardinorafael/pkg/logger"
	"github.com/go-chi/chi"
)

type controller struct {
	ctx       context.Context
	log       logger.Logger
	svc       ServiceInterface
	secretKey string
}

func NewController(
	ctx context.Context,
	log logger.Logger,
	svc ServiceInterface,
	secretKey string,
) *controller {
	return &controller{
		ctx:       ctx,
		log:       log,
		svc:       svc,
		secretKey: secretKey,
	}
}

func (c controller) RegisterRoute(r *chi.Mux) {
	m := middleware.NewWithAuth(c.ctx, c.log, c.secretKey)

	r.Route("/api/v1/organizations", func(r chi.Router) {
		r.Use(m.WithAuth)

		r.Post("/", c.CreateOrg)
		r.Get("/{orgId}", c.GetOrgByID)
		r.Get("/slug/{slug}", c.GetOrgBySlug)
	})
}

func (c controller) CreateOrg(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name    string `json:"name"`
		OwnerID string `json:"owner_id"`
	}

	if err := util.ReadRequestBody(w, r, &body); err != nil {
		if err, ok := err.(ApplicationError); ok {
			NewHttpError(w, err)
			return
		}
		NewHttpError(w, NewInternalServerError(err))
		return
	}

	if err := c.svc.CreateOrg(c.ctx, body.Name, body.OwnerID); err != nil {
		if err, ok := err.(ApplicationError); ok {
			NewHttpError(w, err)
			return
		}
		NewHttpError(w, NewInternalServerError(err))
		return
	}

	util.WriteSuccessResponse(w, http.StatusOK)
}

func (c controller) GetOrgByID(w http.ResponseWriter, r *http.Request) {
	org, err := c.svc.GetOrgByID(c.ctx, chi.URLParam(r, "orgId"))
	if err != nil {
		if err, ok := err.(ApplicationError); ok {
			NewHttpError(w, err)
			return
		}
		NewHttpError(w, NewInternalServerError(err))
		return
	}

	util.WriteJSONResponse(w, http.StatusOK, map[string]any{"org": org})
}

func (c controller) GetOrgBySlug(w http.ResponseWriter, r *http.Request) {
	org, err := c.svc.GetOrgBySlug(c.ctx, chi.URLParam(r, "slug"))
	if err != nil {
		if err, ok := err.(ApplicationError); ok {
			NewHttpError(w, err)
			return
		}
		NewHttpError(w, NewInternalServerError(err))
		return
	}

	util.WriteJSONResponse(w, http.StatusOK, map[string]any{"org": org})
}
