package permission

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

	r.Route("/api/v1/permissions", func(r chi.Router) {
		r.Use(m.WithAuth)

		r.Post("/", c.createPermission)
		r.Get("/{teamId}", c.getAllPermissions)
		r.Delete("/{teamId}/{permissionId}", c.deletePermission)
	})
}

func (c controller) deletePermission(w http.ResponseWriter, r *http.Request) {
	err := c.svc.Delete(
		c.ctx,
		chi.URLParam(r, "teamId"),
		chi.URLParam(r, "permissionId"),
	)
	if err != nil {
		errors.NewHttpError(w, errors.NewBadRequestError(
			"failed to delete permission",
			err,
		))
		return
	}

	WriteSuccessResponse(w, http.StatusOK)
}

func (c controller) getAllPermissions(w http.ResponseWriter, r *http.Request) {
	var p PermissionSearchParams
	var teamId = chi.URLParam(r, "teamId")

	p.Query = ReadQueryString(r.URL.Query(), "q", "")
	p.Limit = ReadQueryInt(r.URL.Query(), "limit", 15)
	p.Page = ReadQueryInt(r.URL.Query(), "page", 1)
	p.Sort = ReadQueryString(r.URL.Query(), "sort", "created")

	res, err := c.svc.GetAll(c.ctx, teamId, p)
	if err != nil {
		errors.NewHttpError(w, errors.NewBadRequestError(
			"failed to get all permissions",
			err,
		))
		return
	}

	WriteJSONResponse(w, http.StatusOK, map[string]any{
		"data": res.Data,
		"meta": res.Meta,
	})
}

func (c controller) createPermission(w http.ResponseWriter, r *http.Request) {
	var body CreatePermissionParams

	if err := ReadRequestBody(w, r, &body); err != nil {
		errors.NewHttpError(w, errors.NewBadRequestError(
			"failed to create permission",
			err,
		))
		return
	}

	if err := c.svc.Create(c.ctx, body); err != nil {
		errors.NewHttpError(w, errors.NewBadRequestError(
			"failed to create permission",
			err,
		))
		return
	}

	WriteSuccessResponse(w, http.StatusOK)
}
