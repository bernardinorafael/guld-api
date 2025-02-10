package role

import (
	"context"
	"net/http"

	"github.com/bernardinorafael/internal/_shared/errors"
	. "github.com/bernardinorafael/internal/_shared/util"

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

		r.Post("/", c.createRole)
		r.Post("/{roleId}", c.createRolePermission)
		r.Get("/{teamId}", c.getAllRoles)
		r.Get("/{teamId}/{roleId}", c.getRole)
	})
}

func (c controller) createRolePermission(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Permissions []string `json:"permissions"`
	}

	if err := ReadRequestBody(w, r, &body); err != nil {
		errors.NewHttpError(w, errors.NewBadRequestError(
			"failed to create role permission",
			err,
		))
		return
	}

	err := c.svc.MakeRolePermissions(
		c.ctx,
		chi.URLParam(r, "roleId"),
		body.Permissions,
	)
	if err != nil {
		errors.NewHttpError(w, errors.NewBadRequestError(
			"failed to create role permission",
			err,
		))
		return
	}

	WriteSuccessResponse(w, http.StatusOK)
}

func (c controller) getRole(w http.ResponseWriter, r *http.Request) {
	role, err := c.svc.GetByID(
		c.ctx,
		chi.URLParam(r, "teamId"),
		chi.URLParam(r, "roleId"),
	)
	if err != nil {
		errors.NewHttpError(w, errors.NewBadRequestError(
			"failed to get role",
			err,
		))
		return
	}

	WriteJSONResponse(w, http.StatusOK, map[string]any{
		"role": role,
	})
}

func (c controller) getAllRoles(w http.ResponseWriter, r *http.Request) {
	var teamId = chi.URLParam(r, "teamId")

	roles, err := c.svc.GetAll(c.ctx, teamId)
	if err != nil {
		errors.NewHttpError(w, errors.NewBadRequestError(
			"failed to get all roles",
			err,
		))
		return
	}

	WriteJSONResponse(w, http.StatusOK, map[string]any{
		"roles": roles,
	})
}

func (c controller) createRole(w http.ResponseWriter, r *http.Request) {
	var body CreateRoleProps

	if err := ReadRequestBody(w, r, &body); err != nil {
		errors.NewHttpError(w, errors.NewBadRequestError(
			"failed to create role",
			err,
		))
		return
	}

	if err := c.svc.Create(c.ctx, body); err != nil {
		errors.NewHttpError(w, errors.NewBadRequestError(
			"failed to create role",
			err,
		))
		return
	}

	WriteSuccessResponse(w, http.StatusOK)
}
