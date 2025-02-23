package team

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

	r.Route("/api/v1/teams", func(r chi.Router) {
		r.Use(m.WithAuth)

		r.Post("/", c.create)
		r.Get("/organization/{orgId}/owner/{ownerId}", c.getAll)
		r.Get("/{teamId}/organization/{orgId}", c.getByID)
		r.Get("/{slug}/organization/{orgId}", c.getBySlug)
		r.Post("/{teamId}/members", c.addMember)

		r.Get("/member/{userId}/organization/{orgId}", c.getByMember)
	})
}

func (c controller) getByMember(w http.ResponseWriter, r *http.Request) {
	team, err := c.svc.GetByMember(
		c.ctx,
		chi.URLParam(r, "orgId"),
		chi.URLParam(r, "userId"),
	)
	if err != nil {
		NewHttpError(w, err)
		return
	}

	util.WriteJSONResponse(w, http.StatusOK, map[string]any{
		"team": team,
	})
}

func (c controller) addMember(w http.ResponseWriter, r *http.Request) {
	var body AddMemberParams
	body.TeamID = chi.URLParam(r, "teamId")

	if err := util.ReadRequestBody(w, r, &body); err != nil {
		NewHttpError(w, err)
		return
	}

	if err := c.svc.AddMember(c.ctx, body); err != nil {
		NewHttpError(w, err)
		return
	}

	util.WriteSuccessResponse(w, http.StatusOK)
}

func (c controller) getByID(w http.ResponseWriter, r *http.Request) {
	team, err := c.svc.GetBySlug(
		c.ctx,
		chi.URLParam(r, "orgId"),
		chi.URLParam(r, "teamId"),
	)
	if err != nil {
		NewHttpError(w, err)
		return
	}

	util.WriteJSONResponse(w, http.StatusOK, map[string]any{
		"team": team,
	})
}

func (c controller) getBySlug(w http.ResponseWriter, r *http.Request) {
	team, err := c.svc.GetBySlug(
		c.ctx,
		chi.URLParam(r, "orgId"),
		chi.URLParam(r, "slug"),
	)
	if err != nil {
		NewHttpError(w, err)
		return
	}

	util.WriteJSONResponse(w, http.StatusOK, map[string]any{
		"team": team,
	})
}

func (c controller) getAll(w http.ResponseWriter, r *http.Request) {
	teams, err := c.svc.GetAll(
		c.ctx,
		chi.URLParam(r, "orgId"),
		chi.URLParam(r, "ownerId"),
	)
	if err != nil {
		NewHttpError(w, err)
		return
	}

	util.WriteJSONResponse(w, http.StatusOK, map[string]any{
		"teams": teams,
	})

}

func (c controller) create(w http.ResponseWriter, r *http.Request) {
	var body CreateTeamParams

	if err := util.ReadRequestBody(w, r, &body); err != nil {
		NewHttpError(w, err)
		return
	}

	if err := c.svc.Create(c.ctx, body); err != nil {
		NewHttpError(w, err)
		return
	}

	util.WriteSuccessResponse(w, http.StatusOK)
}
