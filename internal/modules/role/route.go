package role

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

	r.Route("/api/v1/roles", func(r chi.Router) {
		r.Use(m.WithAuth)

		r.Post("/org/{orgId}", c.createRole)
		r.Get("/org/{orgId}", c.getRoles)
	})
}

func (c controller) getRoles(w http.ResponseWriter, r *http.Request) {
	var input RoleSearchParams

	input.Query = util.ReadQueryString(r.URL.Query(), "q", "")
	input.Limit = util.ReadQueryInt(r.URL.Query(), "limit", 25)
	input.Page = util.ReadQueryInt(r.URL.Query(), "page", 1)
	input.Sort = util.ReadQueryString(r.URL.Query(), "sort", "name")

	res, err := c.svc.FindAll(c.ctx, chi.URLParam(r, "orgId"), input)
	if err != nil {
		NewHttpError(w, err)
		return
	}

	util.WriteJSONResponse(w, http.StatusOK, map[string]any{
		"data": res.Data,
		"meta": res.Meta,
	})
}

func (c controller) createRole(w http.ResponseWriter, r *http.Request) {
	var input CreateRoleProps
	input.OrgID = chi.URLParam(r, "orgId")

	err := util.ReadRequestBody(w, r, &input)
	if err != nil {
		NewHttpError(w, err)
		return
	}

	err = c.svc.Create(c.ctx, input)
	if err != nil {
		NewHttpError(w, err)
		return
	}

	util.WriteSuccessResponse(w, http.StatusCreated)
}
