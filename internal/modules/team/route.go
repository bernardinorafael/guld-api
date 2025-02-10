package team

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

	r.Route("/api/v1/teams", func(r chi.Router) {
		r.Use(m.WithAuth)

		r.Post("/", c.createTeam)
		r.Get("/{ownerId}", c.getAllTeams)
	})
}

func (c controller) getAllTeams(w http.ResponseWriter, r *http.Request) {

	teams, err := c.svc.GetAll(c.ctx, chi.URLParam(r, "ownerId"))
	if err != nil {
		errors.NewHttpError(w, errors.NewBadRequestError(
			"failed to delete email",
			err,
		))
		return
	}

	WriteJSONResponse(w, http.StatusOK, map[string]any{
		"teams": teams,
	})

}

func (c controller) createTeam(w http.ResponseWriter, r *http.Request) {
	var body CreateTeamParams

	if err := ReadRequestBody(w, r, &body); err != nil {
		errors.NewHttpError(w, errors.NewBadRequestError(
			"failed to create team",
			err,
		))
		return
	}

	if err := c.svc.Create(c.ctx, body); err != nil {
		errors.NewHttpError(w, errors.NewBadRequestError(
			"failed to create team",
			err,
		))
		return
	}

	WriteSuccessResponse(w, http.StatusOK)
}
