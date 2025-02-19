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

		r.Post("/org/{orgId}", c.createPermission)
		r.Get("/org/{orgId}", c.getAllPermissions)
		r.Delete("/org/{orgId}/perm/{permId}", c.deletePermission)
		r.Put("/org/{orgId}/perm/{permId}", c.updatePermission)
	})
}

func (c controller) deletePermission(w http.ResponseWriter, r *http.Request) {
	err := c.svc.Delete(
		c.ctx,
		chi.URLParam(r, "orgId"),
		chi.URLParam(r, "permId"),
	)
	if err != nil {
		if err, ok := err.(errors.ApplicationError); ok {
			errors.NewHttpError(w, err)
			return
		}
		errors.NewHttpError(w, errors.NewInternalServerError(err))
		return
	}

	WriteSuccessResponse(w, http.StatusOK)
}

func (c controller) updatePermission(w http.ResponseWriter, r *http.Request) {
	var body UpdatePermissionParams
	body.OrgID = chi.URLParam(r, "orgId")
	body.ID = chi.URLParam(r, "permId")

	if err := ReadRequestBody(w, r, &body); err != nil {
		if err, ok := err.(errors.ApplicationError); ok {
			errors.NewHttpError(w, err)
			return
		}
		errors.NewHttpError(w, errors.NewInternalServerError(err))
		return
	}

	if err := c.svc.Update(c.ctx, body); err != nil {
		if err, ok := err.(errors.ApplicationError); ok {
			errors.NewHttpError(w, err)
			return
		}
	}

	WriteSuccessResponse(w, http.StatusOK)
}

func (c controller) getAllPermissions(w http.ResponseWriter, r *http.Request) {
	var body PermissionSearchParams

	body.Query = ReadQueryString(r.URL.Query(), "q", "")
	body.Limit = ReadQueryInt(r.URL.Query(), "limit", 15)
	body.Page = ReadQueryInt(r.URL.Query(), "page", 1)
	body.Sort = ReadQueryString(r.URL.Query(), "sort", "created")

	res, err := c.svc.GetAll(c.ctx, chi.URLParam(r, "orgId"), body)
	if err != nil {
		if err, ok := err.(errors.ApplicationError); ok {
			errors.NewHttpError(w, err)
			return
		}
		errors.NewHttpError(w, errors.NewInternalServerError(err))
		return
	}

	WriteJSONResponse(w, http.StatusOK, map[string]any{
		"data": res.Data,
		"meta": res.Meta,
	})
}

func (c controller) createPermission(w http.ResponseWriter, r *http.Request) {
	var body CreatePermissionParams
	body.OrgID = chi.URLParam(r, "orgId")

	if err := ReadRequestBody(w, r, &body); err != nil {
		if err, ok := err.(errors.ApplicationError); ok {
			errors.NewHttpError(w, err)
			return
		}
		errors.NewHttpError(w, errors.NewInternalServerError(err))
		return
	}

	if err := c.svc.Create(c.ctx, body); err != nil {
		if err, ok := err.(errors.ApplicationError); ok {
			errors.NewHttpError(w, err)
			return
		}
		errors.NewHttpError(w, errors.NewInternalServerError(err))
		return
	}

	WriteSuccessResponse(w, http.StatusOK)
}
