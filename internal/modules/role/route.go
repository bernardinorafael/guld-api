package role

import (
	"context"
	"net/http"

	"github.com/bernardinorafael/internal/_shared/dto"
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
		r.Get("/org/{orgId}/role/{roleId}", c.getRole)
		r.Delete("/org/{orgId}/role/{roleId}", c.deleteRole)
		r.Patch("/org/{orgId}/role/{roleId}", c.updateRole)

		// Permissions
		r.Patch("/permissions/org/{orgId}/role/{roleId}", c.managePermissions)
	})
}

func (c controller) updateRole(w http.ResponseWriter, r *http.Request) {
	roleId := chi.URLParam(r, "roleId")
	orgId := chi.URLParam(r, "orgId")

	var input UpdateRoleDTO

	err := util.ReadRequestBody(w, r, &input)
	if err != nil {
		NewHttpError(w, err)
		return
	}

	err = c.svc.UpdateRoleInformation(
		c.ctx,
		orgId,
		roleId,
		input,
	)
	if err != nil {
		NewHttpError(w, err)
		return
	}

	util.WriteSuccessResponse(w, http.StatusNoContent)
}

func (c controller) deleteRole(w http.ResponseWriter, r *http.Request) {
	roleId := chi.URLParam(r, "roleId")
	orgId := chi.URLParam(r, "orgId")

	err := c.svc.Delete(c.ctx, orgId, roleId)
	if err != nil {
		NewHttpError(w, err)
		return
	}

	util.WriteSuccessResponse(w, http.StatusNoContent)
}

func (c controller) managePermissions(w http.ResponseWriter, r *http.Request) {
	roleId := chi.URLParam(r, "roleId")

	var input ManagePermissionsProps

	err := util.ReadRequestBody(w, r, &input)
	if err != nil {
		NewHttpError(w, err)
		return
	}

	err = c.svc.ManagePermissions(c.ctx, roleId, input.Permissions)
	if err != nil {
		NewHttpError(w, err)
		return
	}

	util.WriteSuccessResponse(w, http.StatusCreated)
}

func (c controller) getRole(w http.ResponseWriter, r *http.Request) {
	roleId := chi.URLParam(r, "roleId")
	orgId := chi.URLParam(r, "orgId")

	role, err := c.svc.GetRole(c.ctx, orgId, roleId)
	if err != nil {
		NewHttpError(w, err)
		return
	}

	util.WriteJSONResponse(w, http.StatusOK, map[string]any{
		"role": role,
	})
}

func (c controller) getRoles(w http.ResponseWriter, r *http.Request) {
	var query = r.URL.Query()
	var input dto.SearchParams

	input.Query = util.ReadQueryString(query, "q", "")
	input.Limit = util.ReadQueryInt(query, "limit", 25)
	input.Page = util.ReadQueryInt(query, "page", 1)
	input.Sort = util.ReadQueryString(query, "sort", "name")

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
