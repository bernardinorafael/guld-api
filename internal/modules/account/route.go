package account

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
	return &controller{ctx, log, svc, secretKey}
}

func (c controller) RegisterRoute(r *chi.Mux) {
	m := middleware.NewWithAuth(c.ctx, c.log, c.secretKey)

	r.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/register", c.register)
		r.Post("/login", c.login)
		r.Post("/activate/{userId}", c.activate)
	})

	r.Route("/api/v1/accounts/me", func(r chi.Router) {
		r.Use(m.WithAuth)
		r.Get("/", c.getSigned)
	})
}

func (c controller) activate(w http.ResponseWriter, r *http.Request) {
	err := c.svc.ActivateAccount(r.Context(), chi.URLParam(r, "userId"))
	if err != nil {
		if err, ok := err.(ApplicationError); ok {
			NewHttpError(w, err)
			return
		}
		NewHttpError(w, NewInternalServerError(err))
		return
	}

	util.WriteSuccessResponse(w, http.StatusOK)
}

func (c controller) getSigned(w http.ResponseWriter, r *http.Request) {
	acc, err := c.svc.GetSignedInAccount(r.Context())
	if err != nil {
		if err, ok := err.(ApplicationError); ok {
			NewHttpError(w, err)
			return
		}
		NewHttpError(w, NewInternalServerError(err))
		return
	}

	util.WriteJSONResponse(w, http.StatusOK, map[string]any{
		"account": acc,
	})
}

func (c controller) login(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := util.ReadRequestBody(w, r, &body); err != nil {
		if err, ok := err.(ApplicationError); ok {
			NewHttpError(w, err)
			return
		}
		NewHttpError(w, NewInternalServerError(err))
		return
	}

	token, claims, err := c.svc.Login(r.Context(), body.Username, body.Password)
	if err != nil {
		if err, ok := err.(ApplicationError); ok {
			NewHttpError(w, err)
			return
		}
		NewHttpError(w, NewInternalServerError(err))
		return
	}

	payload := AccountPayload{
		AccessToken: token,
		AccountID:   claims.AccountID,
		UserID:      claims.UserID,
		OrgID:       claims.OrgID,
		Username:    claims.Username,
		Email:       claims.Email,
		IssuedAt:    claims.RegisteredClaims.IssuedAt.Unix(),
		ExpiresAt:   claims.RegisteredClaims.ExpiresAt.Unix(),
	}

	util.WriteJSONResponse(w, http.StatusAccepted, payload)
}

func (c controller) register(w http.ResponseWriter, r *http.Request) {
	var body CreateAccountParams

	if err := util.ReadRequestBody(w, r, &body); err != nil {
		if err, ok := err.(ApplicationError); ok {
			NewHttpError(w, err)
			return
		}
		NewHttpError(w, NewInternalServerError(err))
		return
	}

	token, claims, err := c.svc.Register(r.Context(), body)
	if err != nil {
		if err, ok := err.(ApplicationError); ok {
			NewHttpError(w, err)
			return
		}
		NewHttpError(w, NewInternalServerError(err))
		return
	}

	payload := AccountPayload{
		AccessToken: token,
		AccountID:   claims.AccountID,
		UserID:      claims.UserID,
		OrgID:       claims.OrgID,
		Username:    claims.Username,
		Email:       claims.Email,
		IssuedAt:    claims.RegisteredClaims.IssuedAt.Unix(),
		ExpiresAt:   claims.RegisteredClaims.ExpiresAt.Unix(),
	}

	util.WriteJSONResponse(w, http.StatusAccepted, payload)
}
