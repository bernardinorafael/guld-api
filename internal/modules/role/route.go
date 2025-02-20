package role

import (
	"context"
	"net/http"

	"github.com/bernardinorafael/internal/_shared/errors"
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

	r.Route("/api/v1/roles", func(r chi.Router) {
		r.Use(m.WithAuth)

		r.Post("/org/{orgId}", c.createRole)
		r.Get("/org/{orgId}", c.getRoles)
	})
}

func (c controller) getRoles(w http.ResponseWriter, r *http.Request) {
	roles, err := c.svc.FindAll(c.ctx, chi.URLParam(r, "orgId"))
	if err != nil {
		if err, ok := err.(errors.ApplicationError); ok {
			errors.NewHttpError(w, err)
			return
		}
		errors.NewHttpError(w, errors.NewInternalServerError(err))
		return
	}

	util.PrintJSON(roles)

	util.WriteJSONResponse(w, http.StatusOK, map[string]any{
		"roles": roles,
	})
}

func (c controller) createRole(w http.ResponseWriter, r *http.Request) {
	var body CreateRoleProps
	body.OrgID = chi.URLParam(r, "orgId")

	if err := util.ReadRequestBody(w, r, &body); err != nil {
		if err, ok := err.(errors.ApplicationError); ok {
			errors.NewHttpError(w, err)
			return
		}
		errors.NewHttpError(w, errors.NewInternalServerError(err))
		return
	}

	err := c.svc.Create(c.ctx, body)
	if err != nil {
		if err, ok := err.(errors.ApplicationError); ok {
			errors.NewHttpError(w, err)
			return
		}
		errors.NewHttpError(w, errors.NewInternalServerError(err))
		return
	}

	util.WriteSuccessResponse(w, http.StatusCreated)
}
